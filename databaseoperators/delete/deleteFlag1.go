package delete

import (
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"time"
)

type DeleteFlag1 struct {
	request *DeleteJSON
}

type DeleteJSON struct {
	StoneID   datatypes.JSONString `json:"stoneID"`
	StartTime datatypes.JSONInt64        `json:"startTime"`
	EndTime   datatypes.JSONInt64         `json:"endTime"`
}

type (
	ResponseFlag1 struct {
		Succeed struct {
			KwH uint32 `json:"kiloWattHour"`
		} `json:"succeed"`
	}
)

func parseFlag1(message []byte, indexOfMessage int) (*DeleteFlag1, datatypes.Error) {

	requestJSON := &DeleteJSON{}
	request := &DeleteFlag1{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, datatypes.NoError
}

func (requestJSON *DeleteFlag1) marshalBytes(message []byte, indexOfMessage int) datatypes.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
	if err != nil {
		return datatypes.UnMarshallError
	}
	return datatypes.NoError
}

func (requestJSON *DeleteFlag1) checkParameters() datatypes.Error {
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

func (requestJSON *DeleteFlag1) Execute() ([]byte, datatypes.Error) {
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

func (requestJSON *DeleteFlag1) executeDatabase() (*ResponseFlag1, datatypes.Error) {
	var queryTimePart string
	var timeValues []interface{}
	timeValues = append(timeValues, requestJSON.request.StoneID.Value)

	var values string
	var error datatypes.Error
	values = util.UnitkWh
	values = "DELETE " + values + " FROM kwh_by_id_and_time_v2 WHERE id = ? AND time = ?"
	numberOfDeletions, error := requestJSON.request.selectAndInsert("SELECT time FROM w_and_pf_by_id_and_time_v2 WHERE id = ?"+queryTimePart, values, timeValues)
	if !error.IsNull() {
		return nil, error
	}
	response := &ResponseFlag1{}
	response.Succeed.KwH = numberOfDeletions
	return response,  datatypes.NoError
}

func (requestJSON *DeleteJSON) selectAndInsert(selectQuery string, insertQuery string, values []interface{}) (uint32, datatypes.Error) {
	var error datatypes.Error
	iterator, error := cassandra.Query(selectQuery, values...)
	if !error.IsNull() {
		return 0, error
	}
	var timeOfRow time.Time
	var timeArray []time.Time
	for iterator.Scan(&timeOfRow) {
		timeArray = append(timeArray, timeOfRow)
	}

	if err := iterator.Close(); err != nil {
		error = datatypes.CassandraIterator
		error.Message = err.Error()
		return 0, error

	}
	batch, error := cassandra.CreateBatch()
	if !error.IsNull() {
		return 0, error
	}
	for _, valueTime := range timeArray {
		error = cassandra.AddQueryToBatchAndExecuteWhenBatchMax(batch, insertQuery, values[0], valueTime)
		if !error.IsNull() {
			return 0, error
		}
	}

	error = cassandra.ExecuteBatch(batch)
	if !error.IsNull() {
		return 0, error
	}

	return uint32(len(timeArray)), datatypes.NoError
}
