package delete

import (
	"GoCQLTimeSeries/model"
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

func parseFlag2(message []byte, indexOfMessage int) (*DeleteFlag2, model.Error) {

	requestJSON := &DeleteJSON{}
	request := &DeleteFlag2{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, model.NoError
}

func (requestJSON *DeleteFlag2) marshalBytes(message []byte, indexOfMessage int) model.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}
	return model.NoError
}

func (requestJSON *DeleteFlag2) checkParameters() model.Error {
	if !requestJSON.request.StoneID.Valid {
		return model.MissingStoneID
	}

	if !requestJSON.request.StartTime.Valid {
		if !requestJSON.request.EndTime.Valid {
			return model.MissingStartAndEndTime
		}
		return model.MissingStartTime
	}

	if !requestJSON.request.EndTime.Valid {
		return model.MissingEndTime
	}

	return model.NoError
}

func (requestJSON *DeleteFlag2) Execute() ([]byte, model.Error) {
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

func (requestJSON *DeleteFlag2) executeDatabase()(*ResponseFlag2, model.Error) {
	var queryTimePart string
	var timeValues []interface{}
	timeValues = append(timeValues, requestJSON.request.StoneID.Value)
	var values string
	var error model.Error
	values = util.UnitW + ", " + util.Unitpf
	values = "DELETE "+values +" FROM w_and_pf_by_id_and_time_v2 WHERE id = ? AND time = ?"
	numberOfDeletions, error := requestJSON.request.selectAndInsert("SELECT time FROM w_and_pf_by_id_and_time_v2 WHERE id = ?"+queryTimePart, values, timeValues)
	if !error.IsNull() {
		return nil, error
	}
	response := &ResponseFlag2{}
	response.Succeed.WattPowerFactor = numberOfDeletions
	return response,  model.NoError
}