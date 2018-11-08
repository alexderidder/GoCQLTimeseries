package parser

import (
	"CrownstoneServer/server/database/connector"
	"encoding/binary"
	"encoding/json"
	"github.com/gocql/gocql"
	"log"
	"time"
)

type requestResponseJSON struct {
	StoneIDs  []gocql.UUID `json:"stoneIDs"`
	Types     []string     `json:"types"`
	StartTime time.Time    `json:"startTime"`
	EndTime   time.Time    `json:"endTime"`
	Interval  uint32       `json:"interval"`
}

type responseSelectJSON struct {
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

type selectCommand struct {
	message    []byte
}

func (s selectCommand) parseFlag () []byte {
	flag := byteToInt(s.message, 0)

	switch flag {
	case 1:

		return s.parseJSON(s.message[4:])

	default:
		return ParseError(10, "Server doesn't recognise flag")
	}

}

func (s selectCommand) parseJSON(message []byte) []byte {
	request := requestResponseJSON{}
	err := json.Unmarshal(message, &request)
	if err != nil {
		return ParseError(100, "JSON layout is wrong")
	}
	errBytes, requestedTypes := s.checkRequiredParameters(request)
	if errBytes != nil {
		return errBytes
	}
	wAndPfTableQuery, kwhTableQuery := s.createQuery(request, requestedTypes)
	response := responseSelectJSON{request.StartTime, request.EndTime, request.Interval, []stone{}}
	var timeValues []interface{}
	wAndPfTableQuery += " WHERE id = ?"
	kwhTableQuery += " WHERE id = ?"

	if !request.StartTime.IsZero() && !request.EndTime.IsZero() {
		wAndPfTableQuery += ` AND time >= ? AND time <= ? `
		kwhTableQuery += ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, request.StartTime)
		timeValues = append(timeValues, request.EndTime)
	}
	var queryValues []interface{}
	for _, stoneID := range request.StoneIDs {
		queryValues = append([]interface{}{}, stoneID)
		queryValues = append(queryValues, timeValues...)

		stone := stone{}
		stone.StoneID = stoneID
		stone.Fields = []field{}
		var iterator *gocql.Iter

		switch requestedTypes[0] {
		case UnitWAndpf:
			iterator = cassandra.Session.Query(wAndPfTableQuery,  queryValues...).Iter()
			var timeOfRow time.Time
			var w, pf *float32
			var wList, pfList []Data
			for iterator.Scan(&timeOfRow, &w, &pf) {
				if w != nil {
					wList = append(wList, Data{timeOfRow, *w})
				}
				if pf != nil {
					pfList = append(pfList, Data{timeOfRow, *pf})
				}


			}
			if err := iterator.Close(); err != nil {
				log.Fatal(err)
				//TODO: research if error code is needed
			}

			stone.Fields = append(stone.Fields, field{UnitW, wList})
			stone.Fields = append(stone.Fields, field{Unitpf, pfList})
		case UnitW:
			iterator = cassandra.Session.Query(wAndPfTableQuery,  queryValues...).Iter()
			measurements, err := s.iterateStreamWithOneFloat32PerRow(iterator, request.Interval)
			if err != nil {
				return ParseError(300, err.Error())
			}
			stone.Fields = append(stone.Fields, field{UnitW, measurements})
		case Unitpf:
			iterator = cassandra.Session.Query(wAndPfTableQuery,  queryValues...).Iter()
			measurements, err := s.iterateStreamWithOneFloat32PerRow(iterator, request.Interval)
			if err != nil {
				return ParseError(300, err.Error())
			}

			stone.Fields = append(stone.Fields, field{Unitpf, measurements})
		}

		switch requestedTypes[1] {
		case UnitkWh:
			iterator = cassandra.Session.Query(kwhTableQuery, queryValues...).Iter()
			measurements, err := s.iterateStreamWithOneFloat32PerRow(iterator, request.Interval)
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

func  (s selectCommand) iterateStreamWithOneFloat32PerRow(iterator *gocql.Iter, interval uint32) ([]Data, error) {
	var value *float32
	var measurementList []Data
	var timeOfRow time.Time

	for iterator.Scan(&timeOfRow, &value) {
		if value != nil{
			measurementList = append(measurementList, Data{timeOfRow, *value})
		}
	}

	if err := iterator.Close(); err != nil {
		return measurementList, err
		//TODO: error code is needed
	}
	return measurementList, nil
}

func (s selectCommand)  checkRequiredParameters(request requestResponseJSON) ([]byte, [2]string) {

	if len(request.StoneIDs) == 0 {
		return ParseError(100, "StoneID is missing"),[2]string{}
	}
	if len(request.Types) == 0 {
		return ParseError(100, "Type is missing"), [2]string{}
	}
	unknownTypes, test := checkUnknownAndDuplicatedTypes(request.Types)

	if !unknownTypes {
		return ParseError(100, "Unknown types added"), [2]string{}
	}
	return nil, test
}

func (s selectCommand) createQuery(request requestResponseJSON, typePerQuery [2]string) (string, string) {
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

	return kwhAndPfTable, wattTable

}

