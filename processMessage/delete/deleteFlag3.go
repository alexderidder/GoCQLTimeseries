package delete

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
)

type DeleteFlag3 struct {
	request *DeleteOnlyStone
}

type DeleteOnlyStone struct {
	StoneIDs []model.JSONString `json:"stoneID"`
}

type ResponseFlag3 struct {
	Succeed struct {
		WattPowerFactor uint32 `json:"wattAndPowerFactor"`
		KwH             uint32 `json:"kiloWattHour"`
	} `json:"succeed"`
}

func parseFlag3(message []byte, indexOfMessage int) (*DeleteFlag3, model.Error) {
	requestJSON := &DeleteOnlyStone{}
	request := &DeleteFlag3{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, model.NoError
}

func (requestJSON *DeleteFlag3) marshalBytes(message []byte, indexOfMessage int) model.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}
	return model.NoError
}

func (requestJSON *DeleteFlag3) checkParameters() model.Error {
	if len(requestJSON.request.StoneIDs) == 0 {
		return model.MissingStoneID
	}
	return model.NoError
}

func (requestJSON *DeleteFlag3) Execute() ([]byte, model.Error) {
	error := requestJSON.executeDatabase()
	if !error.IsNull() {
		return nil, error
	}
	return append(util.Uint32ToByteArray(2)), model.NoError
}

func (requestJSON *DeleteFlag3) executeDatabase() (model.Error) {
	for _, stoneID := range requestJSON.request.StoneIDs {
		error := cassandra.ExecQuery("DELETE FROM w_and_pf_by_id_and_time_v2 WHERE id = ?", stoneID.Value);
		error = cassandra.ExecQuery("DELETE FROM kwh_by_id_and_time_v2 WHERE id = ?", stoneID.Value);
		if !error.IsNull() {
			return error
		}
	}
	return model.NoError
}
