package _select

import (
	"GoCQLTimeSeries/model"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

type RequestFlag1 struct {
	request *RequestJSON
}
type StoneIDsWithBucketsWithDataPoints struct {
	StoneID               string
	BucketsWithDataPoints []BucketWithDataPoints
}

type BucketWithDataPoints struct {
	Bucket           int64
	EnergyDataPoints []Data
}

func parseFlag1(message []byte, indexOfMessage int) (*RequestFlag1, model.Error) {
	requestJSON := &RequestJSON{}
	request := &RequestFlag1{requestJSON}
	if err := request.marshalBytes(message, indexOfMessage); !err.IsNull() {
		return nil, err
	}

	if err := request.checkParameters(); !err.IsNull() {
		return nil, err
	}

	return request, model.NoError
}

func (requestJSON *RequestFlag1) marshalBytes(message []byte, indexOfMessage int) model.Error {
	err := json.Unmarshal(message[indexOfMessage:], requestJSON.request)
	if err != nil {
		error := model.UnMarshallError
		error.Message = err.Error()
		return error
	}

	return model.NoError
}

func (requestJSON *RequestFlag1) checkParameters() model.Error {
	if len(requestJSON.request.StoneIDs) == 0 {
		return model.MissingStoneID
	}
	return model.NoError
}

func (requestJSON *RequestFlag1) Execute() ([]byte, model.Error) {
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

func (requestJSON *RequestFlag1) executeDatabase() (*ResponseJSON, model.Error) {
	response := ResponseJSON{map[string][]Data{}}
	var partDataList = &[]Data{}
	var dataList = []Data{}
	var list = []StoneIDsWithBucketsWithDataPoints{}
	var timeQuery string
	var timeBuckets []int64
	var timeValues []interface{}

	//Calculate interval for BucketWithDataPoints retrieval
	if requestJSON.request.StartTime.Set && requestJSON.request.EndTime.Set {
		timeBuckets = util.GetIntervalsWithAggr1BetweenUnixStamps(requestJSON.request.StartTime.Value, requestJSON.request.EndTime.Value)
		timeValues = append([]interface{}{}, requestJSON.request.StartTime.Value)
		timeValues = append(timeValues, requestJSON.request.EndTime.Value)
		timeQuery = ` AND time >= ? AND time <= ? `
	} else {
		timeBuckets = util.GetIntervalsWithAggr1BetweenDefaultStartStampAndNow()
		timeValues = append([]interface{}{},util.DefaultYearSince )
		timeValues = append(timeValues, time.Now().UnixNano() / int64(time.Millisecond))
		timeQuery = ` AND time >= ? AND time <= ? `
	}

	for _, stoneID := range requestJSON.request.StoneIDs {
	//Request BucketWithDataPoints data from server
		for _, timeBucket := range timeBuckets {

			queryValues := append([]interface{}{}, stoneID.Value, timeBucket)
			queryValues = append(queryValues, timeValues...)

			iterator, error := cassandra.Query("SELECT time, "+util.UnitkWh+" FROM kWh_by_id_and_time_in_layer2 WHERE id = ? and time_bucket = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			 getAlreadyDownSampledData(iterator, partDataList)
		}
	//Calculate interval for Raw retrieval
		var timeBucketsWeek []int64
		var startTime int64
		var endTime int64
		if length := len(*partDataList)-1; length > -1 {
			startTime = (*partDataList)[length].Time
			endTime = requestJSON.request.EndTime.Value
		} else {
			startTime = util.DefaultDaySince
			endTime = time.Now().UnixNano() / int64(time.Millisecond)
		}
		timeBucketsWeek = util.WeeksBetweenDates(startTime, endTime)
		timeValues = append([]interface{}{}, startTime)
		timeValues = append(timeValues, endTime)

		lastYearBucket := int64(0)

		//Request raw data
		aggrData := []BucketWithDataPoints{}
		for _, timeBucket := range timeBucketsWeek {
			queryValues := append([]interface{}{}, stoneID.Value, timeBucket)
			queryValues = append(queryValues, timeValues...)
			if lastYearBucket != timeBucket/52 {
				aggrData = append(aggrData, BucketWithDataPoints{lastYearBucket, *partDataList} )
				dataList = append(dataList, *partDataList...)
				partDataList = &[]Data{}
				lastYearBucket = timeBucket / 52
			}

			iterator, error := cassandra.Query("SELECT time, "+util.UnitkWh+" FROM kWh_by_id_and_time_in_layer1 WHERE id = ? and time_bucket = ?"+timeQuery, queryValues...)
			if !error.IsNull() {
				return nil, error
			}
			error = downsSampleRawDataToAggr1(partDataList, iterator)
			if !error.IsNull() {
				return nil, error
			}

		}
		aggrData = append(aggrData, BucketWithDataPoints{lastYearBucket, *partDataList} )
		dataList = append(dataList, *partDataList...)
		partDataList = &[]Data{}
		response.Stones[stoneID.Value] = dataList
		list = append(list, StoneIDsWithBucketsWithDataPoints{stoneID.Value, aggrData})

	}
	go doSomething(list)
	return &response, model.NoError

}

func getAlreadyDownSampledData(iterator *gocql.Iter, dataList *[]Data ) {
	var kWh *float64
	var timeOfRow *int64
		for iterator.Scan(&timeOfRow, &kWh) {
		if timeOfRow != nil {
			data := Data{Time: *timeOfRow}
			if *kWh != 0 {
				data.Value.KWH = *kWh
				*dataList = append(*dataList, data)
			}
		}
	}
}

func downsSampleRawDataToAggr1(dataList *[]Data, iterator *gocql.Iter) (model.Error) {
	var y1LastValue *float64
	var x1LastTimestamp *int64
	var y0PreviousValue float64
	var x0PreviousTimestamp int64

	var timeStampEvery5Minutes int64
	var checkPreviousValue bool
	if length := len(*dataList) -1; length > -1 {
		x0PreviousTimestamp = (*dataList)[length].Time
		y0PreviousValue = (*dataList)[length].Value.KWH
		checkPreviousValue = false
	} else{
		checkPreviousValue = true
	}

	for iterator.Scan(&x1LastTimestamp, &y1LastValue) {

		if x1LastTimestamp != nil {
			if checkPreviousValue {
				y0PreviousValue = *y1LastValue
				x0PreviousTimestamp = *x1LastTimestamp
				checkPreviousValue = false
				continue
			}
			// Initialize timeStampEvery5Minutes to rounded value of x0PreviousTimestamp at 5 minutes
			timeStampEvery5Minutes = (x0PreviousTimestamp/1000/60/5 + 1) * 5 * 60 * 1000
			for timeStampEvery5Minutes <= *x1LastTimestamp {
				// If range between 2 data-points > 10 minutes. Store checkPreviousValue data-point
				if *x1LastTimestamp-x0PreviousTimestamp > 1000*60*10 {
					data := Data{Time: x0PreviousTimestamp}
					data.Value.KWH = y0PreviousValue;
					*dataList = append(*dataList, data)
					break
				}

				//Linear interpolation
				kwh := util.LinearInterpolation(x0PreviousTimestamp,*x1LastTimestamp, y0PreviousValue,  *y1LastValue, timeStampEvery5Minutes)
				data := Data{Time: timeStampEvery5Minutes}
				data.Value.KWH = kwh;
				*dataList = append(*dataList, data)
				timeStampEvery5Minutes += 1000 * 60 * 5
			}
		}
		y0PreviousValue = *y1LastValue
		x0PreviousTimestamp = *x1LastTimestamp
	}

	if err := iterator.Close(); err != nil {
		error := model.CassandraIterator
		error.Message = err.Error()
		return error

	}

	return model.NoError

}
func doSomething(dataPerStoneIDs []StoneIDsWithBucketsWithDataPoints){
	batch, error := cassandra.CreateBatch()
	if !error.IsNull() {
		fmt.Println(error)
		return
	}

	for _,dataWithStoneID := range dataPerStoneIDs {
		for _, dataWithBucket := range dataWithStoneID.BucketsWithDataPoints {

			for _, dataPoint := range dataWithBucket.EnergyDataPoints {
				cassandra.AddQueryToBatchAndExecuteBatchTillSuccess(batch, "INSERT INTO kWh_by_id_and_time_in_layer2 (id, time_bucket, time, kwh) VALUES (?,?, ?, ?)", dataWithStoneID.StoneID , dataWithBucket.Bucket, dataPoint.Time, dataPoint.Value.KWH)
			}
			error = cassandra.ExecuteBatch(batch)
			if !error.IsNull() {
				fmt.Println(error)
				return
			}

		}
		error = cassandra.ExecuteBatch(batch)
		if !error.IsNull() {
			fmt.Println(error)
			return


	}
	error = cassandra.ExecuteBatch(batch)
	if !error.IsNull() {
		return
	}

}



}

//cassandra.AddQueryToBatchAndExecuteBatchTillSuccess(batch, "INSERT INTO kWh_by_id_and_time_in_layer2 (id, time_bucket, time, kwh) VALUES (?,?, ?, ?)", stoneID, currentYear, time, kWh)
