package parser

import (
	"CrownstoneServer/server/database/connector"
	"encoding/binary"
	"encoding/json"
	"github.com/gocql/gocql"
	"log"
	"strings"
	"time"
)

type query struct {
	StoneIDs  []gocql.UUID `json:"stoneIDs"`
	Types     []string     `json:"types"`
	StartTime time.Time    `json:"startTime"`
	EndTime   time.Time    `json:"endTime"`
	Interval  uint32       `json:"interval"`
}

type response struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Interval  uint32    `json:"interval"`
	Stones    []stones    `json:"stones"`
}

type stones struct {
	StoneID gocql.UUID `json:"stone_id"`
	Types  []types `json:"types"`
}

type types  struct {
	Type         string    `json:"type"`
	Measurements []float32 `json:"measurements"`
}

var posTypes = []string{"pf", "kwh", "w"}

var format = "2006-01-02 15:04:05"

func parseSelect(message []byte) []byte {
	var request query
	err := json.Unmarshal(message, request)
	if err != nil {
		return ParseError(100, "JSON layout is wrong")
	}
	errBytes, requestedTypes := checkRequiredParameters(request)
	if errBytes != nil {
		return errBytes
	}
	leftQuery, rightQuery := createQueryPerStoneID(request, requestedTypes)

	response := response{request.StartTime, request.EndTime, request.Interval, []stones{}}


	for _, value := range request.StoneIDs {
		stones := stones{}
		stones.StoneID = value
		stones.Types = []types{}
		var iterator *gocql.Iter
		if requestedTypes[0] && requestedTypes[1] {
			iterator = cassandra.Session.Query(leftQuery, value).Iter()
			var value, value2 float32
			var listValue, listValue2 []float32
			for iterator.Scan(&value, &value2) {
				listValue = append(listValue, value)
				listValue2 = append(listValue2, value2)
			}
			if err := iterator.Close(); err != nil {
				log.Fatal(err)
				//TODO: research if error code is needed
			}

			stones.Types = append(stones.Types, types{posTypes[0] , listValue})
			stones.Types = append(stones.Types, types{posTypes[1] , listValue2})
		} else if requestedTypes[0] {
			iterator = cassandra.Session.Query(leftQuery, value).Iter()
			stones.Types = append(stones.Types, types{posTypes[0] , test12(iterator)})
		} else if requestedTypes[1] {
			iterator = cassandra.Session.Query(leftQuery, value).Iter()
			stones.Types = append(stones.Types, types{posTypes[1] , test12(iterator)})

		}

		if requestedTypes[2] {
			iterator = cassandra.Session.Query(rightQuery, value).Iter()
			stones.Types = append(stones.Types, types{posTypes[2] , test12(iterator)})
		}
		response.Stones = append(response.Stones, stones)
	}

	result, err := json.Marshal(response)
	opCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(opCode, 1)
	return append(opCode, result...)
}

//func test1(iterator *gocql.Iter ) []float32{
//	return test12(iterator, 1)[0]
//}
//
//func test12(iterator *gocql.Iter, noOfValue uint16  ) [][]float32 {
//	value := make([]float32, noOfValue)
//	listValue := make([][]float32, noOfValue)
//
//	for iterator.Scan(&value) {
//		listValue = append(listValue, value)
//		//TODO: ADD INTERVAL
//	}
//	if err := iterator.Close(); err != nil {
//		log.Fatal(err)
//		//TODO: research if error code is needed
//	}
//	return listValue
//}



func test12(iterator *gocql.Iter  ) []float32 {
	var value float32
	var listValue []float32

	for iterator.Scan(&value) {
		listValue = append(listValue, value)
		//TODO: ADD INTERVAL
	}
	if err := iterator.Close(); err != nil {
		log.Fatal(err)
		//TODO: research if error code is needed
	}
	return listValue
}


func checkRequiredParameters(request query) ([]byte, []bool) {

	if len(request.StoneIDs) == 0 {
		return ParseError(100, "StoneID is missing"), nil
	}
	if len(request.Types) == 0 {
		return ParseError(100, "Type is missing"), nil
	}
	unknownTypes, test := checkUnknownTypes(request.Types)

	if !unknownTypes {
		return ParseError(100, "Unknown types added"), nil
	}
	return nil, test
}

func createQueryPerStoneID(request query, test []bool) (string, string) {
	whereClause := getWhereClause(request)
	var selectClause = "SELECT "
	var leftQuery string
	if test[0] && test[1] {
		leftQuery = selectClause + "sum(" + posTypes[0] + "), sum(" + posTypes[1] + ")?"
	} else if test[0] {
		leftQuery = selectClause + "sum(" + posTypes[0] + ")"
	} else if test[1] {
		leftQuery = selectClause + "sum(" + posTypes[1] + ")"

	}

	var rightQuery string
	if test[2] {
		rightQuery = selectClause + "sum(" + posTypes[2] + ")"
	}
	leftQuery += whereClause
	rightQuery += whereClause

	return leftQuery, rightQuery

}

func getWhereClause(request query) string {
	whereClause := " FROM mytable WHERE stone_id = ?"

	if !request.StartTime.IsZero() && !request.EndTime.IsZero() {
		whereClause += " AND time >= " + request.StartTime.Format(format) + " AND time <= " + request.EndTime.Format(format)
	}
	return whereClause
}

func checkUnknownTypes(request []string) (bool, []bool) {
	var test = []bool{false, false, false}
	for _, v := range request {
		switch strings.ToLower(v) {
		case posTypes[0]:
			test[0] = true
		case posTypes[1]:
			test[1] = true
		case posTypes[2]:
			test[2] = true
		default:
			return false, nil
		}

	}
	return true, test

}
