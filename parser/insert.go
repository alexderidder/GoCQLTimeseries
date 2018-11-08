package parser

import (
	"CrownstoneServer/server/config"
	"CrownstoneServer/server/database/connector"
	"encoding/binary"
	"encoding/json"
	"github.com/gocql/gocql"
	"time"
)

type insertJSON struct {
	StoneID gocql.UUID `json:"stoneID"`
	Data    []struct {
		Time        time.Time   `json:"time"`
		Watt        JSONFloat32 `json:"watt"`
		PowerFactor JSONFloat32 `json:"pf"`
		KWH         JSONFloat32 `json:"kWh"`
	} `json:"data"`
}

type JSONFloat32 struct {
	Value float32
	Valid bool
	Set   bool
}

type insert struct {
	message []byte
}

var batchSize = int(config.Config.Database.BatchSize)

func (i *JSONFloat32) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true

	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}

	// The key isn't set to null
	var temp float32
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}

func (i insert) parseFlag() []byte {
	flag := byteToInt(i.message, 0)

	switch flag {
	case 1:
		return i.parseFlag1ByteToJSON(i.message[4:])

	default:
		return ParseError(10, "Server doesn't recognise flag")
	}

}

func (i insert) parseFlag1ByteToJSON(message []byte) []byte {
	request := insertJSON{}

	err := json.Unmarshal(message, &request)
	if err != nil {
		return ParseError(100, "JSON layout is wrong")
	}

	errBytes := i.checkParameters(request.StoneID, len(request.Data))
	if errBytes != nil {
		return errBytes
	}
	batch := cassandra.Session.NewBatch(gocql.LoggedBatch)
	batch2 := cassandra.Session.NewBatch(gocql.LoggedBatch)

	//TODO: index isnt constant anymore
	for _, data := range request.Data {

		if data.Watt.Valid {
			if data.PowerFactor.Valid {
				batch.Query("INSERT INTO w_and_pw_by_id_and_time (id, time, w, pf) VALUES (?, ?, ?, ?)", request.StoneID, data.Time, data.Watt.Value, data.PowerFactor.Value)
			} else {
				batch.Query("INSERT INTO w_and_pw_by_id_and_time (id, time, w) VALUES (?, ?, ?)", request.StoneID, data.Time, data.Watt.Value)
			}
		} else if data.PowerFactor.Valid {
			batch.Query("INSERT INTO w_and_pw_by_id_and_time (id, time, pf) VALUES (?, ?, ?)", request.StoneID, data.Time, data.PowerFactor.Value)
		} else {

		}

		if batch.Size()%batchSize == 0 {
			err := cassandra.Session.ExecuteBatch(batch)
			if err != nil {
				return ParseError(100, err.Error())
			}
			batch = cassandra.Session.NewBatch(gocql.LoggedBatch)
		}

		if data.KWH.Valid {
			batch2.Query("INSERT INTO kwh_by_id_and_time (id, time, kwh) VALUES (?, ?, ?)", request.StoneID, data.Time, data.KWH.Value)
		} else {
			continue
		}

		if batch2.Size()%batchSize == 0 {
			err := cassandra.Session.ExecuteBatch(batch2)
			if err != nil {
				return ParseError(100, err.Error())
			}
			batch2 = cassandra.Session.NewBatch(gocql.LoggedBatch)
		}
	}
	err = cassandra.Session.ExecuteBatch(batch)
	if err != nil {
		return ParseError(100, err.Error())
	}

	opCode := make([]byte, 4)
	binary.LittleEndian.PutUint32(opCode, 2)
	return opCode
}

func (insert) checkParameters(stoneID gocql.UUID, dataLength int) ([]byte) {

	if len(stoneID) == 0 {
		return ParseError(100, "StoneID is missing")
	}
	if dataLength == 0 {
		return ParseError(100, "No data")
	}

	return nil
}
