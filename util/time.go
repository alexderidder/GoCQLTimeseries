package util

import (
	"GoCQLTimeSeries/config"
	"GoCQLTimeSeries/datatypes"
	"time"
)

func WeeksBetweenDates(startTime, endTime int64) ([]int64) {
	var timeBucketsWeeks []int64
	var tempTime = startTime
	for tempTime < endTime {
		timeBucketsWeeks = append(timeBucketsWeeks, tempTime/config.RAW_DATA_PARTITION_SPLIT_IN_TIMEFRAME)
		tempTime += config.RAW_DATA_PARTITION_SPLIT_IN_TIMEFRAME
	}
	return timeBucketsWeeks
}

func GetIntervalsWithAggr1BetweenUnixStamps(startTime, endTime int64) ([]int64) {
	var timeBucketsWeeks []int64
	var tempTime = startTime
	for tempTime < endTime {
		timeBucketsWeeks = append(timeBucketsWeeks, tempTime/config.LAYER_1_PARTITION_SPLIT_IN_TIMEFRAME)
		tempTime += config.LAYER_1_PARTITION_SPLIT_IN_TIMEFRAME
	}
	return timeBucketsWeeks
}

func GetIntervalsWithAggr1BetweenDefaultStartStampAndNow() ([]int64) {
	return GetIntervalsWithAggr1BetweenUnixStamps(config.BeginDayForRaw, time.Now().UnixNano()/int64(time.Millisecond))
}

func CalculateBucketsForAggregatedRetrieval(partDataList *[]datatypes.Data, startTime int64, endTime int64) ([]int64, []interface{}) {
	var timeBucketsWeek []int64
	if length := len(*partDataList) - 1; length > -1 {
		startTime = (*partDataList)[length].Time
		timeBucketsWeek = WeeksBetweenDates(startTime, endTime)
		timeValues := append([]interface{}{}, startTime)
		timeValues = append(timeValues, endTime)
		return timeBucketsWeek, timeValues
	} else {
		///return error
		return nil, nil
	}


}

func CalculateBucketsForRawRetrieval(startTime int64, endTime int64) ([]int64, []interface{}) {
	var timeBuckets []int64
	var timeValues []interface{}
		timeBuckets = GetIntervalsWithAggr1BetweenUnixStamps(startTime, endTime)
		timeValues = append([]interface{}{}, startTime)
		timeValues = append(timeValues, endTime)

	return timeBuckets, timeValues
}
