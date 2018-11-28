package _select

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"github.com/gocql/gocql"
	"time"
)

type RequestFlag1 struct {
	request *RequestJSON
}

func parseFlag1(message *[]byte) (*RequestFlag1, model.Error) {
	requestJSON := &RequestJSON{}
	request := &RequestFlag1{requestJSON}
	if err := request.marshalBytes(message); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, model.NoError
}

func (requestJSON *RequestFlag1) marshalBytes(message *[]byte) model.Error {
	err := json.Unmarshal(*message, requestJSON.request)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}

	return model.NoError
}

func (requestJSON *RequestFlag1) checkParameters() model.Error {
	if len(requestJSON.request.StoneIDs) == 0 {
		return model.MissingStoneID
	}
	return model.NoError
}

func (requestJSON *RequestFlag1) Execute() ([]byte, model.Error) {
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

func (requestJSON *RequestFlag1) executeDatabase() (*ResponseJSON, model.Error) {
	var error model.Error
	response := ResponseJSON{requestJSON.request.StartTime, requestJSON.request.EndTime, requestJSON.request.Interval, []Stone{}}

	var timeValues []interface{}
	var timeQuery string
	if !requestJSON.request.StartTime.IsZero() && !requestJSON.request.EndTime.IsZero() {
		timeQuery = ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, requestJSON.request.StartTime)
		timeValues = append(timeValues, requestJSON.request.EndTime)
	}
	var queryValues []interface{}

	for _, stoneID := range requestJSON.request.StoneIDs {
		queryValues = append([]interface{}{}, stoneID.Value)
		queryValues = append(queryValues, timeValues...)

		stone := Stone{}
		stone.StoneID = stoneID.Value

		var iterator *gocql.Iter

		iterator, error = cassandra.Query("SELECT time, "+util.UnitkWh+" FROM kwh_by_id_and_time_v2 WHERE id = ?"+timeQuery, queryValues...)
		if !error.IsNull() {
			return nil, error
		}
		var kWh *float64
		var dataList []Data
		var timeOfRow *time.Time

		for iterator.Scan(&timeOfRow, &kWh) {
			if (timeOfRow != nil) {
				data := Data{Time: *timeOfRow,}
				if kWh != nil {
					data.Value.KWH = *kWh
					dataList = append(dataList, data)
				}
			}
		}
		if err := iterator.Close(); err != nil {
			error = model.CassandraIterator
			error.Message = err.Error()
			return nil, error

		}
		stone.Data = dataList
		response.Stones = append(response.Stones, stone)
	}
	return &response, model.NoError

}
