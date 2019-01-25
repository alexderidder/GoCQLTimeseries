package _select

import (
	"GoCQLTimeSeries/config"
	"GoCQLTimeSeries/databaseoperators/datatypes/energy"
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
)

type RequestFlag1 struct {
	request *RequestJSON
}

func parseFlag1(message []byte, indexOfMessage int) (*RequestFlag1, datatypes.Error) {
	requestJSON := &RequestJSON{}
	request := &RequestFlag1{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}
	return request, datatypes.NoError
}

func (requestJSON *RequestFlag1) marshalBytes(message []byte, indexOfMessage int) datatypes.Error {
	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
	if err != nil {
		return datatypes.UnMarshallError
	}

	return datatypes.NoError
}

func (requestJSON *RequestFlag1) checkParameters() datatypes.Error {
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

func (requestJSON *RequestFlag1) Execute() ([]byte, datatypes.Error) {
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

func (requestJSON *RequestFlag1) executeDatabase() (*ResponseJSON, datatypes.Error) {
	response := ResponseJSON{map[string][]datatypes.Data{}}
	var partDataList []datatypes.Data
	var dataList []datatypes.Data
	var layer1DataForInsertion []datatypes.StoneIDsWithBucketsWithDataPoints

	beginTime := requestJSON.request.StartTime.Value -config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	endTime := requestJSON.request.EndTime.Value + config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	//Calculate interval for Raw retrieval
	timeBuckets, timeValues := util.CalculateBucketsForRawRetrieval(beginTime, endTime)
	timeQuery := ` AND time >= ? AND time <= ? `

	for _, stoneID := range requestJSON.request.StoneIDs {
		//Request bucket data from server
		for _, timeBucket := range timeBuckets {

			queryValues := append([]interface{}{}, stoneID.Value, timeBucket)
			queryValues = append(queryValues, timeValues...)

			iterator, error := cassandra.Query("SELECT time, "+util.UnitkWh+" FROM kWh_by_id_and_time_in_layer2 WHERE id = ? and time_bucket = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			energy.GetAlreadyDownSampledData(iterator, &partDataList)
		}

		timeBucketsLayer1, timeValues := util.CalculateBucketsForAggregatedRetrieval(&partDataList, beginTime, endTime)
		lastYearBucket := int64(0)

		//Request raw data
		var aggrData []datatypes.BucketWithDataPoints
		var lastDataPoint = datatypes.Data{Time:-1,}
		for _, timeBucket := range timeBucketsLayer1 {
			queryValues := append([]interface{}{}, stoneID.Value, timeBucket)
			queryValues = append(queryValues, timeValues...)
			if lastYearBucket != timeBucket/52 {
				aggrData = append(aggrData, datatypes.BucketWithDataPoints{lastYearBucket, partDataList})
				dataList = append(dataList, partDataList...)
				partDataList = []datatypes.Data{}
				lastYearBucket = timeBucket / 52
			}

			iterator, error := cassandra.Query("SELECT time, "+util.UnitkWh+" FROM kWh_by_id_and_time_in_layer1 WHERE id = ? and time_bucket = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			lastDataPoint, error = energy.DownsSampleRawDataToAggr1(&partDataList, lastDataPoint, iterator)
			if !error.IsNull() {
				return nil, error
			}

		}
		aggrData = append(aggrData, datatypes.BucketWithDataPoints{lastYearBucket, partDataList})
		dataList = append(dataList, partDataList...)
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
	go energy.InsertAggregatedDataInCassandra(layer1DataForInsertion)
	return &response, datatypes.NoError

}

