package _select

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"github.com/gocql/gocql"
	"time"
)

type RequestFlag2 struct {
	request *RequestJSON
}
type RequestJSON struct {
	StoneIDs  []model.JSONString `json:"stoneIDs"`
	StartTime time.Time          `json:"startTime"`
	EndTime   time.Time          `json:"endTime"`
	Interval  uint32             `json:"interval"`
}

type ResponseJSON struct {
	StartTime time.Time    `json:"startTime"`
	EndTime   time.Time    `json:"endTime"`
	Interval  uint32       `json:"interval"`
	Stones    []Stone `json:"stones"`
}

type Stone struct {
	StoneID string      `json:"stoneID"`
	Data    []Data `json:"data"`
}

type Data struct {
	Time  time.Time `json:"time"`
	Value struct {
		Wattage     float32 `json:"w,omitempty"`
		PowerFactor float32 `json:"pf,omitempty"`
		KWH float64 `json:"kWh,omitempty"`
	} `json:"value"`
}

func parseFlag2(message *[]byte) (*RequestFlag2, model.Error) {
	requestJSON := &RequestJSON{}
	request := &RequestFlag2{requestJSON}
	if err := request.marshalBytes(message); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, model.NoError
}

func (requestJSON *RequestFlag2) marshalBytes(message *[]byte) model.Error {
	err := json.Unmarshal(*message, requestJSON.request)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}

	return model.NoError
}

func (requestJSON *RequestFlag2) checkParameters() model.Error {
	if len(requestJSON.request.StoneIDs) == 0 {
		return model.MissingStoneID
	}
	return model.NoError
}

func (requestJSON *RequestFlag2) Execute() ([]byte, model.Error) {
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

func (requestJSON *RequestFlag2) executeDatabase() (*ResponseJSON, model.Error) {
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
		iterator, error = cassandra.Query("SELECT time, "+util.UnitW+", "+util.Unitpf+" FROM w_and_pf_by_id_and_time_v2 WHERE id = ?"+timeQuery, queryValues...)
		if !error.IsNull() {
			return nil, error
		}
		var dataList []Data
		var timeOfRow *time.Time
		var w, pf *float32

		for iterator.Scan(&timeOfRow, &w, &pf) {
			if (timeOfRow != nil) {
				var data= Data{Time: *timeOfRow,}
				if w != nil {
					data.Value.Wattage = *w
					if pf != nil {
						data.Value.PowerFactor = *pf
					}
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
