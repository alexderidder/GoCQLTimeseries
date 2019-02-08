package insert

import (

	"GoCQLTimeSeries/config"
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"sort"
	"strconv"
	"time"
)

type RequestFlag1 struct {
	StoneID datatypes.JSONString `json:"stoneID"`
	Data    []Data `json:"data"`
}

type Data struct {
	Time  int64 `json:"time"`
	Value struct {
		KWH datatypes.JSONFloat64 `json:"kWh"`
	} `json:"value"`
}

type ResponseFlag1 struct {
	Succeed struct {
		KwH float32 `json:"kWh"`
	}
}

func parseFlag1(message []byte, indexOfMessage int) (*RequestFlag1, datatypes.Error) {
	requestJSON := &RequestFlag1{}

	if err := requestJSON.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := requestJSON.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return requestJSON, datatypes.NoError
}

func (requestJSON *RequestFlag1) marshalBytes(message []byte, indexOfMessage int) datatypes.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON)
	if err != nil {
		return datatypes.UnMarshallError
	}
	return datatypes.NoError
}

func (requestJSON *RequestFlag1) checkParameters() datatypes.Error {
	if !requestJSON.StoneID.Valid {
		return datatypes.MissingStoneID
	}

	if len(requestJSON.Data) == 0 {
		return datatypes.MissingData
	}

	return datatypes.NoError
}

func (requestJSON *RequestFlag1) Execute() ([]byte, datatypes.Error) {

	response, error := requestJSON.ExecuteDatabase()
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

func (requestJSON *RequestFlag1) ExecuteDatabase() (*ResponseFlag1, datatypes.Error) {
	response := &ResponseFlag1{}
	var error datatypes.Error
	batch, error := cassandra.CreateBatch()
	if !error.IsNull() {
		return nil, error
	}

	sort.Slice(requestJSON.Data, func(i, j int) bool {
		return requestJSON.Data[i].Time < requestJSON.Data[j].Time
	})

	if indexOfLastDataPoint := len(requestJSON.Data) - 1; indexOfLastDataPoint > -1 {
		if requestJSON.Data[0].Time < config.BeginDayForRaw {
			return nil, datatypes.Error{100, "Can't insert before " + strconv.FormatInt(config.BeginDayForRaw, 10)}
		} else if nowInMilli := time.Now().UnixNano() / 1000; requestJSON.Data[indexOfLastDataPoint].Time > nowInMilli {
			return nil, datatypes.Error{100, "Can't insert after " + strconv.FormatInt(nowInMilli, 10)}
		}
	}

	//var timestampInMilliSeconds int64
	var currentWeek int64 = 0
	for _, data := range requestJSON.Data {
		if data.Value.KWH.Valid {
			response.Succeed.KwH++
			//Floors automatic, timestamp can't be below ..
			if tempweek := data.Time / config.RAW_DATA_PARTITION_SPLIT_IN_TIMEFRAME; tempweek != currentWeek {
				err := cassandra.ExecuteAndClearBatch(batch)
				if !err.IsNull() {
					return nil, err
				}

				for err = cassandra.ExecQuery("INSERT INTO kwh_inserted_in_layer1 (time_bucket, id) VALUES (?, ?)", tempweek, requestJSON.StoneID.Value); !error.IsNull(); {

				}

				currentWeek = tempweek
			}
			err := cassandra.AddQueryToBatchAndExecuteWhenBatchMax(batch, "INSERT INTO kWh_by_id_and_time_in_layer1 (id, time_bucket, time, kwh) VALUES (?,?, ?, ?)", requestJSON.StoneID.Value, currentWeek, data.Time, data.Value.KWH.Value)
			if !err.IsNull() {
				return nil, err
			}

		}

	}

	error = cassandra.ExecuteBatch(batch)
	if !error.IsNull() {
		return nil, error
	}

	return response, datatypes.NoError
}
