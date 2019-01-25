package _select

import (
	"GoCQLTimeSeries/config"
	"GoCQLTimeSeries/databaseoperators/datatypes/power"
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"fmt"
)

type RequestFlag2 struct {
	request *RequestJSON
}
type RequestJSON struct {
	StoneIDs  []datatypes.JSONString `json:"stoneIDs"`
	StartTime datatypes.JSONInt64    `json:"startTime"`
	EndTime   datatypes.JSONInt64    `json:"endTime"`
	Interval  uint32             `json:"interval"`
}

type ResponseJSON struct {
	//StartTime int64             `json:"startTime"`
	//EndTime   int64             `json:"endTime"`
	//Bucket  uint32            `json:"interval"`
	Stones    map[string][]datatypes.Data `json:"stones"`
}


func parseFlag2(message []byte, indexOfMessage int) (*RequestFlag2, datatypes.Error) {
	requestJSON := &RequestJSON{}
	request := &RequestFlag2{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, datatypes.NoError
}

func (requestJSON *RequestFlag2) marshalBytes(message []byte, indexOfMessage int) datatypes.Error {
	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
	if err != nil {
		return datatypes.UnMarshallError
	}

	return datatypes.NoError
}

func (requestJSON *RequestFlag2) checkParameters() datatypes.Error {
	if len(requestJSON.request.StoneIDs) == 0 {
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

func (requestJSON *RequestFlag2) Execute() ([]byte, datatypes.Error) {
	response, error := requestJSON.executeDatabase()
	if !error.IsNull() {
		return nil, error
	}

	if !requestJSON.request.StartTime.Valid {
		if !requestJSON.request.EndTime.Valid {
			return nil, datatypes.MissingStartAndEndTime
		}
		return nil,  datatypes.MissingStartTime
	}

	responseJSONBytes, err := json.Marshal(response)
	if err != nil {
		error := datatypes.MarshallError
		error.Message = err.Error()
		return nil, error
	}
	return append(util.Uint32ToByteArray(1), responseJSONBytes...), datatypes.NoError
}


func (requestJSON *RequestFlag2) executeDatabase() (*ResponseJSON, datatypes.Error) {
	response := ResponseJSON{map[string][]datatypes.Data{}}
	var partDataList []datatypes.Data
	var dataList []datatypes.Data
	var layer1DataForInsertion []datatypes.StoneIDsWithBucketsWithDataPoints
	//Calculate interval for Raw retrieval
	beginTime := requestJSON.request.StartTime.Value -config.MILLISECONDS_INTERVAL_FOR_POWER_AGGREGATION
	endTime := requestJSON.request.EndTime.Value + config.MILLISECONDS_INTERVAL_FOR_POWER_AGGREGATION
	timeBuckets, timeValues := util.CalculateBucketsForRawRetrieval(beginTime, endTime)
	timeQuery := ` AND time >= ? AND time <= ? `
	fmt.Println(timeBuckets, timeValues)
	for _, stoneID := range requestJSON.request.StoneIDs {
		//Request bucket data from server
		for _, timeBucket := range timeBuckets {

			queryValues := append([]interface{}{}, stoneID.Value, timeBucket)
			queryValues = append(queryValues, timeValues...)

			iterator, error := cassandra.Query("SELECT time, "+util.UnitW+", "+util.Unitpf+" FROM w_and_pf_by_id_and_time_layer1 WHERE id = ? and time_bucket = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			power.GetAlreadyDownSampledData(iterator, &partDataList)
		}

		timeBucketsRaw, timeValues := util.CalculateBucketsForAggregatedRetrieval(&partDataList, beginTime, endTime)
		lastYearBucket := int64(0)
		fmt.Println(timeBucketsRaw, timeValues)
		//Request raw data
		var aggrData []datatypes.BucketWithDataPoints
		var lastSample power.AggregatePower
		for _, timeBucket := range timeBucketsRaw {
			queryValues := append([]interface{}{}, stoneID.Value, timeBucket)
			queryValues = append(queryValues, timeValues...)
			if lastYearBucket != timeBucket/52 {
				aggrData = append(aggrData, datatypes.BucketWithDataPoints{lastYearBucket, partDataList})
				dataList = append(dataList, partDataList...)
				partDataList = []datatypes.Data{}
				lastYearBucket = timeBucket / 52
			}

			iterator, error := cassandra.Query("SELECT time, "+util.UnitW+", "+util.Unitpf+" FROM w_and_pf_by_id_and_time_raw WHERE id = ? and time_bucket = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			lastSample, error = power.DownsSampleRawDataToAggr1(&partDataList, lastSample, iterator)
			if !error.IsNull() {
				return nil, error
			}

		}
		aggrData = append(aggrData, datatypes.BucketWithDataPoints{lastYearBucket, partDataList})
		dataList = append(dataList, partDataList...)
		fmt.Println(len(dataList))
		fmt.Println(len(partDataList))
		partDataList = []datatypes.Data{}
		if length := len(dataList); length > 2 {
			if dataList[0].Time < requestJSON.request.StartTime.Value {
				dataList = dataList[1:]
				length --
			}
			if dataList[length-1].Time > requestJSON.request.StartTime.Value {
				dataList = dataList[:length-1]
			}
		}
		response.Stones[stoneID.Value] = dataList
		dataList = []datatypes.Data{}
		layer1DataForInsertion = append(layer1DataForInsertion, datatypes.StoneIDsWithBucketsWithDataPoints{stoneID.Value, aggrData})
	}
	go power.InsertAggregatedDataInCassandra(layer1DataForInsertion)
	return &response, datatypes.NoError

}


		//
		//var iterator *gocql.Iter
		//iterator, error = cassandra.Query("SELECT time, "+util.UnitW+", "+util.Unitpf+" FROM w_and_pf_by_id_and_time_v2 WHERE id = ?"+timeQuery, queryValues...)
		//if !error.IsNull() {
		//	return nil, error
		//}
		//var dataList []datatypes.Data
		//var timeOfRow *time.Time
		//var w, pf *float32
		//
		//for iterator.Scan(&timeOfRow, &w, &pf) {
		//	if timeOfRow != nil {
		//		var data = datatypes.Data{Time: timeOfRow.Unix(),}
		//		if w != nil {
		//			data.Value.Wattage = *w
		//			if pf != nil {
		//				data.Value.PowerFactor = *pf
		//			}
		//			dataList = append(dataList, data)
		//		}
		//	}
		//}
		//if err := iterator.Close(); err != nil {
		//	error = datatypes.CassandraIterator
		//	error.Message = err.Error()
		//	return nil, error
		//
		//}
