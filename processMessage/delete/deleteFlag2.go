package delete

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/util"
	"encoding/json"
)

type DeleteFlag2 struct {
	request *DeleteJSON
}


func parseFlag2(message *[]byte) (*DeleteFlag2, model.Error) {

	requestJSON := &DeleteJSON{}
	request := &DeleteFlag2{requestJSON}
	if err := request.marshalBytes(message); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, model.NoError
}

func (requestJSON *DeleteFlag2) marshalBytes(message *[]byte) model.Error {

	err := json.Unmarshal(*message, requestJSON.request)
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

	return model.NoError
}

func (requestJSON *DeleteFlag2) Execute() ([]byte, model.Error) {
	err := requestJSON.executeDatabase()
	if !err.IsNull() {
		return nil, err
	}

	return util.Uint32ToByteArray(2), model.NoError
}

func (requestJSON *DeleteFlag2) executeDatabase() model.Error {
	var queryTimePart string
	var timeValues []interface{}
	timeValues = append(timeValues, requestJSON.request.StoneID.Value)
	if !requestJSON.request.StartTime.IsZero() && !requestJSON.request.EndTime.IsZero() {
		queryTimePart = ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, requestJSON.request.StartTime)
		timeValues = append(timeValues, requestJSON.request.EndTime)
	}
	var values string
	var error model.Error
	values = util.UnitW + ", " + util.Unitpf
	values = "DELETE "+values +" FROM w_and_pf_by_id_and_time_v2 WHERE id = ? AND time = ?"
	error = requestJSON.request.selectAndInsert("SELECT time FROM w_and_pf_by_id_and_time_v2 WHERE id = ?"+queryTimePart, values, timeValues)
	if !error.IsNull() {
		return error
	}
	return model.NoError
}