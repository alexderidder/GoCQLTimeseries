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

type RequestFlag2 struct {
	StoneID datatypes.JSONString `json:"stoneID"`
	Data    []struct {
		Time  int64 `json:"time"`
		Value struct {
			Watt        datatypes.JSONFloat32 `json:"watt"`
			PowerFactor datatypes.JSONFloat32 `json:"pf"`
		} `json:"value"`
	} `json:"data"`
}

type ResponseFlag2 struct {
	Succeed struct {
		WattPowerFactor float32 `json:"wPf"`
	}
}

func parseFlag2(message []byte, indexOfMessage int) (*RequestFlag2, datatypes.Error) {
	requestJSON := &RequestFlag2{}
	if err := requestJSON.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := requestJSON.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return requestJSON, datatypes.NoError
}

func (requestJSON *RequestFlag2) marshalBytes(message []byte, indexOfMessage int) datatypes.Error {

	err := json.Unmarshal(message[indexOfMessage:], requestJSON)
	if err != nil {
		return datatypes.UnMarshallError
	}
	return datatypes.NoError
}

func (requestJSON *RequestFlag2) checkParameters() datatypes.Error {
	if !requestJSON.StoneID.Valid {
		return datatypes.MissingStoneID
	}

	if len(requestJSON.Data) == 0 {
		return datatypes.MissingData
	}

	return datatypes.NoError
}

func (requestJSON *RequestFlag2) Execute() ([]byte, datatypes.Error) {
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

func (requestJSON *RequestFlag2) executeDatabase() (*ResponseFlag2, datatypes.Error) {
	response := &ResponseFlag2{}
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
	var timestampInMilliSeconds int64
	var currentWeek int64 = 0
	for _, data := range requestJSON.Data {

		if data.Value.Watt.Valid {
			response.Succeed.WattPowerFactor++
			if !data.Value.PowerFactor.Valid {
				data.Value.PowerFactor.Value = float32(1)
			}

			timestampInMilliSeconds = data.Time
			//Floors automatic, timestamp can't be below ..
			if tempweek := timestampInMilliSeconds / config.RAW_DATA_PARTITION_SPLIT_IN_TIMEFRAME; tempweek != currentWeek {
				err := cassandra.ExecuteAndClearBatch(batch)
				if !err.IsNull() {
					return nil, err
				}

				for err = cassandra.ExecQuery("INSERT INTO w_pf_inserted_in_layer1 (time_bucket, id) VALUES (?, ?)", tempweek, requestJSON.StoneID.Value); !error.IsNull(); {

				}

				currentWeek = tempweek
			}
			//fmt.Println(currentWeek)
			err := cassandra.AddQueryToBatchAndExecuteWhenBatchMax(batch, "INSERT INTO w_and_pf_by_id_and_time_raw  (id, time_bucket, time, w, pf) VALUES (?, ?, ?, ?, ?)", requestJSON.StoneID.Value, currentWeek, timestampInMilliSeconds, data.Value.Watt.Value, data.Value.PowerFactor.Value)
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
