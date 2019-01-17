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
	StartTime model.JSONInt64    `json:"startTime"`
	EndTime   model.JSONInt64    `json:"endTime"`
	Interval  uint32             `json:"interval"`
}

type ResponseJSON struct {
	//StartTime int64             `json:"startTime"`
	//EndTime   int64             `json:"endTime"`
	//Bucket  uint32            `json:"interval"`
	Stones    map[string][]Data `json:"stones"`
}

type Data struct {
	Time  int64 `json:"time"`
	Value struct {
		Wattage     float32 `json:"w,omitempty"`
		PowerFactor float32 `json:"pf,omitempty"`
		KWH         float64 `json:"kWh,omitempty"`
	} `json:"value"`
}

func parseFlag2(message []byte, indexOfMessage int) (*RequestFlag2, model.Error) {
	requestJSON := &RequestJSON{}
	request := &RequestFlag2{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, model.NoError
}

func (requestJSON *RequestFlag2) marshalBytes(message []byte, indexOfMessage int) model.Error {
	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
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
	response := ResponseJSON{ Stones: map[string][]Data{}}

	var timeValues []interface{}
	var timeQuery string
	if requestJSON.request.StartTime.Set && requestJSON.request.EndTime.Set {
		timeQuery = ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, requestJSON.request.StartTime.Value)
		timeValues = append(timeValues, requestJSON.request.EndTime.Value)
	}
	var queryValues []interface{}

	for _, stoneID := range requestJSON.request.StoneIDs {
		queryValues = append([]interface{}{}, stoneID.Value)
		queryValues = append(queryValues, timeValues...)

		var iterator *gocql.Iter
		iterator, error = cassandra.Query("SELECT time, "+util.UnitW+", "+util.Unitpf+" FROM w_and_pf_by_id_and_time_v2 WHERE id = ?"+timeQuery, queryValues...)
		if !error.IsNull() {
			return nil, error
		}
		var dataList []Data
		var timeOfRow *time.Time
		var w, pf *float32

		for iterator.Scan(&timeOfRow, &w, &pf) {
			if timeOfRow != nil {
				var data = Data{Time: timeOfRow.Unix(),}
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
		response.Stones[stoneID.Value] = dataList
	}
	return &response, model.NoError

}
