package parser

import (
	"CrownstoneServer/server/database/connector"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"time"
)

type deleteJSON struct {
	StoneID   gocql.UUID `json:"stoneIDs"`
	Types     []string   `json:"types"`
	StartTime time.Time  `json:"startTime"`
	EndTime   time.Time  `json:"endTime"`
}


type delete struct {
	message    []byte
}

func (d delete) parseFlag() []byte {
	flag := byteToInt(d.message, 0)
	d.message = d.message[4:]
	switch flag {
	case 1:
		return d.parseJSON()
	default:
		return  ParseError(10, "Server doesn't recognise flag")

	}

}

func (d delete) parseJSON() []byte {
	request := deleteJSON{}
	err := json.Unmarshal(d.message, &request)
	if err != nil {
		 ParseError(100, "JSON layout is wrong")
	}
	errBytes, filteredTypes := d.checkRequiredParameters(request)
	if errBytes != nil {
		return errBytes
	}
	createdQueries := d.createQuery(filteredTypes)
	var queryTimePart string
	var timeValues []interface{}
	timeValues = append(timeValues, request.StoneID)
	if !request.StartTime.IsZero() && !request.EndTime.IsZero() {
		fmt.Println("kom hier")
		queryTimePart = ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, request.StartTime)
		timeValues = append(timeValues, request.EndTime)
	}
	if len(createdQueries[0][0]) != 0 {
		 errBytes = d.selectAndInsert(createdQueries[1][0]+queryTimePart, createdQueries[0][0], timeValues)
		if errBytes != nil {
			return errBytes
		}
	}
	if len(createdQueries[1][1]) != 0 {
		errBytes = d.selectAndInsert(createdQueries[1][1]+queryTimePart, createdQueries[0][1], timeValues)
		if errBytes != nil {
			return errBytes
		}
	}

	opCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(opCode, 2)
	return  opCode
}

func (d delete) selectAndInsert(selectQuery string, insertQuery string, values []interface{}) []byte{
	var err error
		iterator := cassandra.Session.Query(selectQuery, values...).Iter()
		var timeOfRow time.Time
		var timeArray []time.Time
		for iterator.Scan(&timeOfRow) {
			timeArray = append(timeArray, timeOfRow)
		}

		if err := iterator.Close(); err != nil {
			log.Fatal(err)
			//TODO: research if error code is needed
		}
		batch := cassandra.Session.NewBatch(gocql.LoggedBatch)
		for index, valueTime := range timeArray {
			batch.Query(insertQuery, values[0], valueTime)

			if index%batchSize == 0 {
				err := cassandra.Session.ExecuteBatch(batch)
				if err != nil {

					return ParseError(100, err.Error())
				}
				batch = cassandra.Session.NewBatch(gocql.LoggedBatch)
			}
		}

		if batch.Size() > 0{
			err = cassandra.Session.ExecuteBatch(batch)
			if err != nil {
				return ParseError(100, err.Error())
			}
		}
	return nil
}

func (d delete) checkRequiredParameters(query2 deleteJSON) ( []byte, [2]string){

	if len(query2.StoneID) == 0 {
		//TODO: CHECK possibly wrong
		return  ParseError(100, "StoneID is missing"), [2]string{}
	}

	if len(query2.Types) == 0 {
		return ParseError(100, "Type is missing"), [2]string{}
	}
	unknownTypes, filteredTypes := checkUnknownAndDuplicatedTypes(query2.Types)

	if !unknownTypes {
		return ParseError(100, "Unknown types added"),[2]string{}

	}
	return nil, filteredTypes

}

func (d delete) createQuery(filteredTypes [2]string)  [2][2]string{
	var createdQueries [2][2]string
	switch filteredTypes[0] {
	case UnitWAndpf:
		createdQueries[0][0] = "UPDATE w_and_pw_by_id_and_time SET w = null, pf = null WHERE id = ? AND time = ?"
		createdQueries[1][0] = "SELECT time FROM w_and_pw_by_id_and_time WHERE id = ?"
	case UnitW:
		createdQueries[0][0] = "UPDATE w_and_pw_by_id_and_time SET w = null WHERE id = ? AND time = ?"
		createdQueries[1][0] = "SELECT time FROM w_and_pw_by_id_and_time WHERE id = ?"
	case Unitpf:
		createdQueries[0][0] = "UPDATE w_and_pw_by_id_and_time SET pf = null WHERE id = ? AND time = ?"
		createdQueries[1][0] = "SELECT time FROM w_and_pw_by_id_and_time WHERE id = ?"

	}
	switch filteredTypes[1] {
	case UnitkWh:
		createdQueries[0][1] = "UPDATE kwh_by_id_and_time SET kwh = null WHERE id = ? AND time = ?"
		createdQueries[1][1] = "SELECT time FROM kwh_by_id_and_time WHERE id = ?"
	}
	return createdQueries

}
