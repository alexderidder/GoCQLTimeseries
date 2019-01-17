package util

import (
	"fmt"
	"time"
)

const msInWeek = 604800000
const weeksInYear = 52
const msInYear = msInWeek * weeksInYear

//Since 2015 January
const DefaultYearSince int64 = 1546300800000
//Since month ago
var DefaultDaySince  =  time.Now().AddDate(0,0,-30).UnixNano() / int64(time.Millisecond)

func WeeksBetweenDates(startTime, endTime int64) ([]int64){
	var timeBucketsWeeks []int64
	var tempTime = startTime
	for tempTime < endTime {
		timeBucketsWeeks = append(timeBucketsWeeks, tempTime/msInWeek)
		tempTime += msInWeek
		fmt.Println(startTime, endTime)
	}
	return timeBucketsWeeks
}

func GetIntervalsWithAggr1BetweenUnixStamps(startTime, endTime int64) ([]int64){
	var timeBucketsWeeks []int64
	var tempTime = startTime
	for tempTime < endTime {
		timeBucketsWeeks = append(timeBucketsWeeks, tempTime/msInYear)
		tempTime += msInYear
	}
	return timeBucketsWeeks
}

func GetIntervalsWithAggr1BetweenDefaultStartStampAndNow() ([]int64){
	return GetIntervalsWithAggr1BetweenUnixStamps(DefaultYearSince, time.Now().UnixNano() / int64(time.Millisecond))
}
