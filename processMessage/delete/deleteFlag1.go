package delete

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"time"
)

type DeleteFlag1 struct {
	request *DeleteJSON
}

type DeleteJSON struct {
	StoneID   model.JSONString `json:"stoneID"`
	StartTime time.Time        `json:"startTime"`
	EndTime   time.Time        `json:"endTime"`
}

func parseFlag1(message *[]byte) (*DeleteFlag1, model.Error) {

	requestJSON := &DeleteJSON{}
	request := &DeleteFlag1{requestJSON}
	if err := request.marshalBytes(message); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, model.NoError
}

func (requestJSON *DeleteFlag1) marshalBytes(message *[]byte) model.Error {

	err := json.Unmarshal(*message, requestJSON.request)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}
	return model.NoError
}

func (requestJSON *DeleteFlag1) checkParameters() model.Error {
	if !requestJSON.request.StoneID.Valid {
		return model.MissingStoneID
	}

	return model.NoError
}

func (requestJSON *DeleteFlag1) Execute() ([]byte, model.Error) {
	err := requestJSON.executeDatabase()
	if !err.IsNull() {
		return nil, err
	}

	return util.Uint32ToByteArray(2), model.NoError
}

func (requestJSON *DeleteFlag1) executeDatabase() model.Error {
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
	values = util.UnitkWh
	values = "DELETE "+values +" FROM kwh_by_id_and_time_v2 WHERE id = ? AND time = ?"
	error = requestJSON.request.selectAndInsert("SELECT time FROM w_and_pf_by_id_and_time_v2 WHERE id = ?"+queryTimePart, values, timeValues)
	if !error.IsNull() {
		return error
	}
	return model.NoError
}

func (requestJSON *DeleteJSON) selectAndInsert(selectQuery string, insertQuery string, values []interface{}) model.Error {
	var error model.Error
	iterator, error := cassandra.Query(selectQuery, values...)
	if !error.IsNull() {
		return error
	}
	var timeOfRow time.Time
	var timeArray []time.Time
	for iterator.Scan(&timeOfRow) {
		timeArray = append(timeArray, timeOfRow)
	}

	if err := iterator.Close(); err != nil {
		error = model.CassandraIterator
		error.Message = err.Error()
		return error

	}
	batch, error := cassandra.CreateBatch()
	if !error.IsNull() {
		return error
	}
	for _, valueTime := range timeArray {
		error = cassandra.AddQueryToBatch(batch, insertQuery, values[0], valueTime)
		if !error.IsNull() {
			return error
		}
	}

	error = cassandra.ExecuteBatch(batch)
	if !error.IsNull() {
		return error
	}

	return model.NoError
}
