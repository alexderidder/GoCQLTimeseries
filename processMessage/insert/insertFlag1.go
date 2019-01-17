package insert

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"sort"
)

type RequestFlag1 struct {
	StoneID model.JSONString `json:"stoneID"`
	Data    []struct {
		Time  int64 `json:"time"`
		Value struct {
			KWH model.JSONFloat64 `json:"kWh"`
		} `json:"value"`
	} `json:"data"`
}

type ResponseFlag1 struct {
	Succeed struct {
		KwH float32 `json:"kWh"`
	}
}

func parseFlag1(message []byte, indexOfMessage int) (*RequestFlag1, model.Error) {
	requestJSON := &RequestFlag1{}


	if err := requestJSON.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := requestJSON.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return requestJSON, model.NoError
}

func (requestJSON *RequestFlag1) marshalBytes(message []byte, indexOfMessage int) model.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}
	return model.NoError
}

func (requestJSON *RequestFlag1) checkParameters() model.Error {
	if !requestJSON.StoneID.Valid {
		return model.MissingStoneID
	}

	if len(requestJSON.Data) == 0 {
		return model.MissingData
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

func (requestJSON *RequestFlag1) executeDatabase() (*ResponseFlag1, model.Error) {
	response := &ResponseFlag1{}
	var error model.Error
	batch, error := cassandra.CreateBatch()
	if !error.IsNull() {
		return nil, error
	}

	sort.Slice(requestJSON.Data, func(i, j int) bool {
		return requestJSON.Data[i].Time < requestJSON.Data[j].Time
	})
	var timestampInSeconds int64
	var currentWeek int64 = 0
	for _, data := range requestJSON.Data {

		if data.Value.KWH.Valid {
			response.Succeed.KwH++
			timestampInSeconds = data.Time
			if tempweek := timestampInSeconds/604800000; tempweek != currentWeek {
				err := cassandra.ExecuteAndClearBatch(batch)
				if !err.IsNull() {
					return nil, err
				}

				for err = cassandra.ExecQuery("INSERT INTO kwh_inserted_in_layer1 (time_bucket, id) VALUES (?, ?)", tempweek, requestJSON.StoneID.Value); !error.IsNull(); {

				}

				currentWeek = tempweek
			}
			//fmt.Println(currentWeek)
			err := cassandra.AddQueryToBatchAndExecuteWhenBatchMax(batch, "INSERT INTO kWh_by_id_and_time_in_layer1 (id, time_bucket, time, kwh) VALUES (?,?, ?, ?)", requestJSON.StoneID.Value, currentWeek, timestampInSeconds, data.Value.KWH.Value)
			if !err.IsNull() {
				return nil, err
			}

		}

	}

	error = cassandra.ExecuteBatch(batch)
	if !error.IsNull() {
		return nil, error
	}

	return response, model.NoError
}
