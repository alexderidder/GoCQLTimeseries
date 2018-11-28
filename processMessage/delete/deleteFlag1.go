package delete

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"log"
	"time"
)

type DeleteJSON struct {
	StoneID   model.JSONUUID `json:"stoneID"`
	Types     []string       `json:"types"`
	StartTime time.Time      `json:"startTime"`
	EndTime   time.Time      `json:"endTime"`
}

func parseFlag1(message *[]byte) (*DeleteJSON, model.Error) {
	requestJSON := &DeleteJSON{}
	if err := requestJSON.marshalBytes(message); !err.IsNull() {
		return nil, err
	}

	if err := requestJSON.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return requestJSON, model.NoError
}

func (requestJSON *DeleteJSON) marshalBytes(message *[]byte) model.Error {

	err := json.Unmarshal(*message, requestJSON)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}
	return model.NoError
}

func (requestJSON *DeleteJSON) checkParameters() model.Error {
	if !requestJSON.StoneID.Valid {
		return model.MissingStoneID
	}

	if len(requestJSON.Types) == 0 {
		return model.MissingType
	}
	requestJSON.Types = util.CheckUnknownAndDuplicatedTypes(requestJSON.Types)

	return model.NoError
}

func (requestJSON *DeleteJSON) Execute() ([]byte, model.Error) {
	err := requestJSON.executeDatabase()
	if !err.IsNull() {
		return nil, err
	}

	return util.Uint32ToByteArray(2), model.NoError
}

func (requestJSON *DeleteJSON) executeDatabase() model.Error {
	var queryTimePart string
	var timeValues []interface{}
	timeValues = append(timeValues, requestJSON.StoneID.Value)
	if !requestJSON.StartTime.IsZero() && !requestJSON.EndTime.IsZero() {
		queryTimePart = ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, requestJSON.StartTime)
		timeValues = append(timeValues, requestJSON.EndTime)
	}

	var values string
	var error model.Error
	switch requestJSON.Types[0] {
	case util.UnitWAndpf:
		values = util.UnitW + ", " + util.Unitpf
	case util.UnitW:
		values = util.UnitW
	case util.Unitpf:
		values = util.Unitpf
	default:
		goto Skip
	}
	values = "DELETE "+values +" FROM w_and_pf_by_id_and_time WHERE id = ? AND time = ?"

	error = requestJSON.selectAndInsert("SELECT time FROM w_and_pf_by_id_and_time WHERE id = ?"+queryTimePart, values, timeValues)
	if !error.IsNull() {
		return error
	}
Skip:
	switch requestJSON.Types[1] {
	case util.UnitkWh:
		error = requestJSON.selectAndInsert("SELECT time FROM kwh_by_id_and_time WHERE id = ?"+queryTimePart, "DELETE "+util.UnitkWh+" FROM kwh_by_id_and_time WHERE id = ? AND time = ?", timeValues)
		if !error.IsNull() {
			return error
		}
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
		log.Fatal(err)
		//TODO: research if error code is needed
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
