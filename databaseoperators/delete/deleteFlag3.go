package delete

import (
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
)

type DeleteFlag3 struct {
	request *DeleteOnlyStone
}

type DeleteOnlyStone struct {
	StoneIDs []datatypes.JSONString `json:"stoneID"`
}

type ResponseFlag3 struct {
	Succeed struct {
		WattPowerFactor uint32 `json:"wattAndPowerFactor"`
		KwH             uint32 `json:"kiloWattHour"`
	} `json:"succeed"`
}

func parseFlag3(message []byte, indexOfMessage int) (*DeleteFlag3, datatypes.Error) {
	requestJSON := &DeleteOnlyStone{}
	request := &DeleteFlag3{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, datatypes.NoError
}

func (requestJSON *DeleteFlag3) marshalBytes(message []byte, indexOfMessage int) datatypes.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
	if err != nil {
		return datatypes.UnMarshallError
	}
	return datatypes.NoError
}

func (requestJSON *DeleteFlag3) checkParameters() datatypes.Error {
	if len(requestJSON.request.StoneIDs) == 0 {
		return datatypes.MissingStoneID
	}

	return datatypes.NoError
}

func (requestJSON *DeleteFlag3) Execute() ([]byte, datatypes.Error) {
	error := requestJSON.executeDatabase()
	if !error.IsNull() {
		return nil, error
	}
	return append(util.Uint32ToByteArray(2)), datatypes.NoError
}

func (requestJSON *DeleteFlag3) executeDatabase() (datatypes.Error) {
	for _, stoneID := range requestJSON.request.StoneIDs {
		error := cassandra.ExecQuery("DELETE FROM w_and_pf_by_id_and_time_v2 WHERE id = ?", stoneID.Value);
		error = cassandra.ExecQuery("DELETE FROM kwh_by_id_and_time_v2 WHERE id = ?", stoneID.Value);
		if !error.IsNull() {
			return error
		}
	}
	return datatypes.NoError
}
