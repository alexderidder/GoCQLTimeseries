package delete

import (
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/util"
	"encoding/json"
)

type DeleteFlag2 struct {
	request *DeleteJSON
}

type ResponseFlag2 struct {
	Succeed struct {
		WattPowerFactor uint32 `json:"wattAndPowerFactor"`
	} `json:"succeed"`
}

func parseFlag2(message []byte, indexOfMessage int) (*DeleteFlag2, datatypes.Error) {

	requestJSON := &DeleteJSON{}
	request := &DeleteFlag2{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, datatypes.NoError
}

func (requestJSON *DeleteFlag2) marshalBytes(message []byte, indexOfMessage int) datatypes.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
	if err != nil {
		return datatypes.UnMarshallError
	}
	return datatypes.NoError
}

func (requestJSON *DeleteFlag2) checkParameters() datatypes.Error {
	if !requestJSON.request.StoneID.Valid {
		return datatypes.MissingStoneID
	}

	if !requestJSON.request.StartTime.Valid {
		if !requestJSON.request.EndTime.Valid {
			return datatypes.MissingStartAndEndTime
		}
		return datatypes.MissingStartTime
	}

	if !requestJSON.request.EndTime.Valid {
		return datatypes.MissingEndTime
	}

	return datatypes.NoError
}

func (requestJSON *DeleteFlag2) Execute() ([]byte, datatypes.Error) {
	response, error := requestJSON.executeDatabase()
	if !error.IsNull() {
		return nil, error
	}
	responseJSONBytes, err := json.Marshal(response)
	if err != nil {
		error := datatypes.MarshallError
		error.Message = err.Error()
		return nil, error
	}

	return append(util.Uint32ToByteArray(1), responseJSONBytes...), datatypes.NoError
}

func (requestJSON *DeleteFlag2) executeDatabase()(*ResponseFlag2, datatypes.Error) {
	var queryTimePart string
	var timeValues []interface{}
	timeValues = append(timeValues, requestJSON.request.StoneID.Value)
	var values string
	var error datatypes.Error
	values = util.UnitW + ", " + util.Unitpf
	values = "DELETE "+values +" FROM w_and_pf_by_id_and_time_v2 WHERE id = ? AND time = ?"
	numberOfDeletions, error := requestJSON.request.selectAndInsert("SELECT time FROM w_and_pf_by_id_and_time_v2 WHERE id = ?"+queryTimePart, values, timeValues)
	if !error.IsNull() {
		return nil, error
	}
	response := &ResponseFlag2{}
	response.Succeed.WattPowerFactor = numberOfDeletions
	return response,  datatypes.NoError
}