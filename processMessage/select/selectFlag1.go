package _select

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"github.com/gocql/gocql"
	"log"
	"time"
)

type RequestJSON struct {
	StoneIDs  []model.JSONUUID `json:"stoneIDs"`
	Types     []string         `json:"types"`
	StartTime time.Time        `json:"startTime"`
	EndTime   time.Time        `json:"endTime"`
	Interval  uint32           `json:"interval"`
}

type ResponseJSON struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Interval  uint32    `json:"interval"`
	Stones    []Stone   `json:"stones"`
}

type Stone struct {
	StoneID gocql.UUID `json:"stoneID"`
	Fields  []Field    `json:"fields"`
}

type Field struct {
	Field string `json:"field"`
	Data  []Data `json:"Data"`
}

type Data struct {
	Time  time.Time `json:"time"`
	Value float32   `json:"value"`
}

func parseFlag1(message *[]byte) (*RequestJSON, model.Error) {
	requestJSON := &RequestJSON{}
	if err := requestJSON.marshalBytes(message); !err.IsNull() {
		return nil, err
	}

	if err := requestJSON.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return requestJSON, model.NoError
}

func (requestJSON *RequestJSON) marshalBytes(message *[]byte) model.Error {
	err := json.Unmarshal(*message, requestJSON)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}

	return model.NoError
}

func (requestJSON *RequestJSON) checkParameters() model.Error {
	if len(requestJSON.StoneIDs) == 0 {
		return model.MissingStoneID
	}

	if len(requestJSON.Types) == 0 {
		return model.MissingType
	}

	requestJSON.Types = util.CheckUnknownAndDuplicatedTypes(requestJSON.Types)

	return model.NoError
}

func (requestJSON *RequestJSON) Execute() ([]byte, model.Error) {
	response, error := requestJSON.executeDatabase()
	if !error.IsNull() {
		return nil, error
	}

	responseJSONBytes, err := json.Marshal(response)
	if err != nil {
		error := model.MarshallError
		error.Message = err.Error()
		return nil, error
	}
	return append(util.Uint32ToByteArray(1), responseJSONBytes...), model.NoError
}

func (requestJSON *RequestJSON) executeDatabase() (*ResponseJSON, model.Error) {
	var error model.Error
	response := ResponseJSON{requestJSON.StartTime, requestJSON.EndTime, requestJSON.Interval, []Stone{}}

	var timeValues []interface{}
	var timeQuery string
	if !requestJSON.StartTime.IsZero() && !requestJSON.EndTime.IsZero() {
		timeQuery = ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, requestJSON.StartTime)
		timeValues = append(timeValues, requestJSON.EndTime)
	}
	var queryValues []interface{}

	for _, stoneID := range requestJSON.StoneIDs {
		queryValues = append([]interface{}{}, stoneID.Value)
		queryValues = append(queryValues, timeValues...)

		stone := Stone{}
		stone.StoneID = stoneID.Value
		stone.Fields = []Field{}
		var iterator *gocql.Iter

		switch requestJSON.Types[0] {
		case util.UnitWAndpf:
			iterator, error = cassandra.Query("SELECT time, "+util.UnitW+", "+util.Unitpf+" FROM w_and_pf_by_id_and_time WHERE id = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
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
			stone.Fields = append(stone.Fields, Field{util.UnitW, wList})
			stone.Fields = append(stone.Fields, Field{util.Unitpf, pfList})
		case util.UnitW:
			iterator, error = cassandra.Query("SELECT time, "+util.UnitW+" FROM w_and_pf_by_id_and_time WHERE id = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			measurements, err := requestJSON.iterateStreamWithOneFloat32PerRow(iterator, requestJSON.Interval)
			if err != nil {
				return &response, model.Error{300, err.Error()}
			}
			stone.Fields = append(stone.Fields, Field{util.UnitW, measurements})
		case util.Unitpf:
			iterator, error = cassandra.Query("SELECT time, "+util.Unitpf+" FROM w_and_pf_by_id_and_time WHERE id = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			measurements, err := requestJSON.iterateStreamWithOneFloat32PerRow(iterator, requestJSON.Interval)
			if err != nil {
				return &response, model.Error{300, err.Error()}
			}

			stone.Fields = append(stone.Fields, Field{util.Unitpf, measurements})
		}

		switch requestJSON.Types[1] {
		case util.UnitkWh:
			iterator, error = cassandra.Query("SELECT time, "+util.UnitkWh+" FROM kwh_by_id_and_time WHERE id = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			measurements, err := requestJSON.iterateStreamWithOneFloat32PerRow(iterator, requestJSON.Interval)
			if err != nil {
				return &response, model.Error{300, err.Error()}
			}
			stone.Fields = append(stone.Fields, Field{util.UnitkWh, measurements})
		}
		response.Stones = append(response.Stones, stone)
	}
	return &response, model.NoError

}

func (requestJSON *RequestJSON) iterateStreamWithOneFloat32PerRow(iterator *gocql.Iter, interval uint32) ([]Data, error) {
	var value *float32
	var measurementList []Data
	var timeOfRow time.Time

	for iterator.Scan(&timeOfRow, &value) {
		if value != nil {
			measurementList = append(measurementList, Data{timeOfRow, *value})
		}
	}

	if err := iterator.Close(); err != nil {
		//fmt.Println(err)
		return measurementList, err
		//TODO: error code is needed
	}

	return measurementList, nil
}
