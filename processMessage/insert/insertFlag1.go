package insert

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"time"
)

type Request struct {
	StoneID model.JSONString `json:"stoneID"`
	Data    []struct {
		Time        time.Time         `json:"time"`
		Watt        model.JSONFloat32 `json:"watt"`
		PowerFactor model.JSONFloat32 `json:"pf"`
		KWH         model.JSONFloat64 `json:"kWh"`
	} `json:"data"`
}

func parseFlag1(message *[]byte) (*Request, model.Error) {
	requestJSON := &Request{}
	if err := requestJSON.marshalBytes(message); !err.IsNull() {
		return nil, err
	}

	if err := requestJSON.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return requestJSON, model.NoError
}

func (requestJSON *Request) marshalBytes(message *[]byte) model.Error {

	err := json.Unmarshal(*message, requestJSON)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}
	return model.NoError
}

func (requestJSON *Request) checkParameters() model.Error {
	if !requestJSON.StoneID.Valid {
		return model.MissingStoneID
	}

	if len(requestJSON.Data) == 0 {
		return model.MissingData
	}

	return model.NoError
}

func (requestJSON *Request) Execute() ([]byte, model.Error) {
	err := requestJSON.executeDatabase()
	if !err.IsNull() {
		return nil, err
	}

	return util.Uint32ToByteArray(2), model.NoError
}

func (requestJSON *Request) executeDatabase() model.Error {
	var error model.Error
	batch, error := cassandra.CreateBatch()
	if !error.IsNull() {
		return error
	}
	batch2, error := cassandra.CreateBatch()
	if !error.IsNull() {
		return error
	}
	for _, data := range requestJSON.Data {
		if data.Watt.Valid {
			if data.PowerFactor.Valid {
				error = cassandra.AddQueryToBatch(batch, "INSERT INTO w_and_pf_by_id_and_time_v2  (id, time, w, pf) VALUES (?, ?, ?, ?)", requestJSON.StoneID.Value, data.Time, data.Watt.Value, data.PowerFactor.Value)
			} else{
				error = cassandra.AddQueryToBatch(batch, "INSERT INTO w_and_pf_by_id_and_time_v2 (id, time, w, pf) VALUES (?, ?, ?, ?)", requestJSON.StoneID.Value, data.Time, data.Watt.Value, float32(1))
			}

			if !error.IsNull() {
				return error
			}
		} else {			// else?
		}

		if data.KWH.Valid {
			err := cassandra.AddQueryToBatch(batch2, "INSERT INTO kwh_by_id_and_time_v2 (id, time, kwh) VALUES (?, ?, ?)", requestJSON.StoneID.Value, data.Time, data.KWH.Value)
			if !err.IsNull() {
				return err
			}
		} else {
			// else?
		}

	}

	error = cassandra.ExecuteBatch(batch)
	if !error.IsNull() {
		return error
	}

	error = cassandra.ExecuteBatch(batch2)
	if !error.IsNull() {
		return error
	}

	return model.NoError
}
