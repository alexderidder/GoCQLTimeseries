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

type responseFormat struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Interval  uint32    `json:"interval"`
	Stones    []stone   `json:"stones"`
}

type stone struct {
	StoneID gocql.UUID `json:"stoneID"`
	Fields  []field    `json:"fields"`
}

type field struct {
	Field string `json:"field"`
	Data  []Data `json:"Data"`
}

type Data struct {
	Time  time.Time `json:"time"`
	Value float32   `json:"value"`
}

const (
	UnitW      string = "w"
	Unitpf     string = "pf"
	UnitWAndpf string = "w_pf"
	UnitkWh    string = "kwh"
)

var format = "2006-01-02 15:04:05"

func parseSelect(message []byte) []byte {
	flag := byteToInt(message, 0)

	switch flag {
	case 1:

		return parseSelectJSON(message[4:])

	default:
		return ParseError(10, "Server doesn't recognise flag")
	}

}

func parseSelectJSON(message []byte) []byte {
	request := query{}
	err := json.Unmarshal(message, &request)
	if err != nil {
		return ParseError(100, "JSON layout is wrong")
	}
	errBytes, requestedTypes := checkRequiredParameters(request)
	if errBytes != nil {
		return errBytes
	}
	wAndPfTableQuery, kwhTableQuery := createQueryPerStoneID(request, requestedTypes)
	response := responseFormat{request.StartTime, request.EndTime, request.Interval, []stone{}}

	for _, stoneID := range request.StoneIDs {
		stone := stone{}
		stone.StoneID = stoneID
		stone.Fields = []field{}
		var iterator *gocql.Iter

		switch requestedTypes[0] {
		case UnitWAndpf:
			iterator = cassandra.Session.Query(wAndPfTableQuery, stoneID).Iter()
			var timeOfRow time.Time
			var w, pf float32
			var kWhList, pfList []Data
			for iterator.Scan(&timeOfRow, &w, &pf) {
				kWhList = append(kWhList, Data{timeOfRow, w})
				pfList = append(pfList, Data{timeOfRow, pf})
			}
			if err := iterator.Close(); err != nil {
				log.Fatal(err)
				//TODO: research if error code is needed
			}

			stone.Fields = append(stone.Fields, field{UnitkWh, kWhList})
			stone.Fields = append(stone.Fields, field{Unitpf, pfList})
		case UnitW:
			iterator = cassandra.Session.Query(wAndPfTableQuery, stoneID).Iter()
			measurements, err := iterateStreamWithOneFloat32PerRow(iterator, request.Interval)
			if err != nil {
				return ParseError(300, err.Error())
			}
			stone.Fields = append(stone.Fields, field{UnitW, measurements})
		case Unitpf:
			iterator = cassandra.Session.Query(wAndPfTableQuery, stoneID).Iter()
			measurements, err := iterateStreamWithOneFloat32PerRow(iterator, request.Interval)
			if err != nil {
				return ParseError(300, err.Error())
			}

			stone.Fields = append(stone.Fields, field{Unitpf, measurements})
		}

		switch requestedTypes[1] {
		case UnitkWh:
			iterator = cassandra.Session.Query(kwhTableQuery, stoneID).Iter()
			measurements, err := iterateStreamWithOneFloat32PerRow(iterator, request.Interval)
			if err != nil {
				return ParseError(300, err.Error())
			}
			stone.Fields = append(stone.Fields, field{UnitkWh, measurements})
		}
		response.Stones = append(response.Stones, stone)
	}

	result, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err.Error())
	}
	opCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(opCode, 1)
	return append(opCode, result...)
}


func iterateStreamWithOneFloat32PerRow(iterator *gocql.Iter, interval uint32) ([]Data, error) {
	var value float32
	var measurementList []Data
	var timeOfRow time.Time

	for iterator.Scan(&timeOfRow, &value) {
		measurementList = append(measurementList, Data{timeOfRow, value})
	}

	if err := iterator.Close(); err != nil {
		return measurementList, err
		//TODO: error code is needed
	}
	return measurementList, nil
}

func checkRequiredParameters(request query) ([]byte, []string) {

	if len(request.StoneIDs) == 0 {
		return ParseError(100, "StoneID is missing"), nil
	}
	if len(request.Types) == 0 {
		return ParseError(100, "Type is missing"), nil
	}
	unknownTypes, test := checkUnknownAndDuplicatedTypes(request.Types)

	if !unknownTypes {
		return ParseError(100, "Unknown types added"), nil
	}
	return nil, test
}

func createQueryPerStoneID(request query, typePerQuery []string) (string, string) {
	whereClause := getWhereClause(request)
	var selectClause = "SELECT time, "
	var kwhAndPfTable string

	switch typePerQuery[0] {
	case UnitWAndpf:
		kwhAndPfTable = selectClause + UnitW + ", " + Unitpf + " FROM w_and_pw_by_id_and_time"
	case UnitW:
		kwhAndPfTable = selectClause + UnitW + " FROM w_and_pw_by_id_and_time"
	case Unitpf:
		kwhAndPfTable = selectClause + Unitpf + " FROM w_and_pw_by_id_and_time"
	}
	var wattTable string
	switch typePerQuery[1] {
	case UnitkWh:
		wattTable = selectClause + UnitkWh + " FROM kwh_by_id_and_time"
	}

	kwhAndPfTable += whereClause
	wattTable += whereClause

	return kwhAndPfTable, wattTable

}

func getWhereClause(request query) string {
	whereClause := " WHERE id = ?"

	if !request.StartTime.IsZero() && !request.EndTime.IsZero() {
		whereClause += " AND time >= " + request.StartTime.Format(format) + " AND time <= " + request.EndTime.Format(format)
	}
	return whereClause
}

func checkUnknownAndDuplicatedTypes(request []string) (bool, []string) {
	var typeList = []bool{false, false, false}
	for _, v := range request {
		switch strings.ToLower(v) {
		case UnitW:
			typeList[0] = true
		case Unitpf:
			typeList[1] = true
		case UnitkWh:
			typeList[2] = true
		default:
			return false, nil
		}

	}
	var typePerQuery = make([]string, 2)
	if typeList[0] && typeList[1] {
		typePerQuery[0] = UnitWAndpf
	} else if typeList[0] {
		typePerQuery[0] = UnitW
	} else if typeList[1] {
		typePerQuery[0] = Unitpf
	}

	if typeList[2] {
		typePerQuery[1] = UnitkWh
	}

	return true, typePerQuery

}
