package insert

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
)

type RequestFlag2 struct {
	StoneID model.JSONString `json:"stoneID"`
	Data    []struct {
		Time  int64 `json:"time"`
		Value struct {
			Watt        model.JSONFloat32 `json:"watt"`
			PowerFactor model.JSONFloat32 `json:"pf"`
		} `json:"value"`
	} `json:"data"`
}

type ResponseFlag2 struct {
	Succeed struct {
		WattPowerFactor float32 `json:"wPf"`
	}
}

func parseFlag2(message []byte, indexOfMessage int) (*RequestFlag2, model.Error) {
	requestJSON := &RequestFlag2{}
	if err := requestJSON.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := requestJSON.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return requestJSON, model.NoError
}

func (requestJSON *RequestFlag2) marshalBytes(message []byte, indexOfMessage int) model.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}
	return model.NoError
}

func (requestJSON *RequestFlag2) checkParameters() model.Error {
	if !requestJSON.StoneID.Valid {
		return model.MissingStoneID
	}

	if len(requestJSON.Data) == 0 {
		return model.MissingData
	}

	return model.NoError
}

func (requestJSON *RequestFlag2) Execute() ([]byte, model.Error) {
	response ,error := requestJSON.executeDatabase()
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

func (requestJSON *RequestFlag2) executeDatabase() (*ResponseFlag2, model.Error) {
	response := &ResponseFlag2{}
	var error model.Error
	batch, error := cassandra.CreateBatch()
	if !error.IsNull() {
		return nil, error
	}
	for _, data := range requestJSON.Data {

		if data.Value.Watt.Valid {
			response.Succeed.WattPowerFactor++
			if data.Value.PowerFactor.Valid {

				error = cassandra.AddQueryToBatchAndExecuteWhenBatchMax(batch, "INSERT INTO w_and_pf_by_id_and_time_v2  (id, time, w, pf) VALUES (?, ?, ?, ?)", requestJSON.StoneID.Value, data.Time, data.Value.Watt.Value, data.Value.PowerFactor.Value)
			} else {
				error = cassandra.AddQueryToBatchAndExecuteWhenBatchMax(batch, "INSERT INTO w_and_pf_by_id_and_time_v2 (id, time, w, pf) VALUES (?, ?, ?, ?)", requestJSON.StoneID.Value, data.Time, data.Value.Watt.Value, float32(1))
			}

			if !error.IsNull() {
				return nil, error
			}
		} else {
		}

	}


	error = cassandra.ExecuteBatch(batch)
	if !error.IsNull() {
		return nil, error
	}

	return response, model.NoError
}
