package energy

import (
	"GoCQLTimeSeries/config"
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/server/cassandra"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type A []struct {
	a int64
	b float64
}

const timestampSinceEpoch = 1548168745 - 1548168745 %config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
func TestWithNoTimestampsFromDatabase(t *testing.T) {

	var m = &cassandra.IteratorMock{}

	m.On("Scan", mock.Anything, mock.Anything).Return(false)

	var aggregatedList []datatypes.Data
	aggregatedList = []datatypes.Data{}
	lastValue := datatypes.Data{Time: -1,}
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, datatypes.Data{Time: -1,}, lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, []datatypes.Data{}, aggregatedList, "Check if no value is added to aggregatedList")

	lastValue = datatypes.Data{}
	lastValue.Time = 100
	lastValue.Value.KWH = 150
	checkValue := lastValue
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, checkValue, lastValue, "Last value is datapoint with time 100 & KWH 150")
	assert.Equal(t, []datatypes.Data{lastValue}, aggregatedList, "Check if lastValue is added to aggregatedList")

	lastValue = datatypes.Data{Time: -1,}
	aggregatedList = []datatypes.Data{{Time: 100,}}
	copyList := aggregatedList
	assert.Equal(t, datatypes.Data{Time: -1,}, lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, copyList, aggregatedList, "Check if no value is added to aggregatedList")

	lastValue = datatypes.Data{}
	lastValue.Time = 100
	lastValue.Value.KWH = 150
	checkValue = lastValue
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, checkValue, lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, []datatypes.Data{{Time: 100,}, lastValue}, aggregatedList, "Check if lastValue is added to aggregatedList")

}

func TestCheckWithPreciseAggregationTimestamps(t *testing.T) {

	var m, _ = &cassandra.IteratorMock{}, fmt.Errorf("e")

	var now int64 = timestampSinceEpoch
	var kwh float64
	dataPointsFromDatabase := make([]datatypes.Data, 10)
	for i := 0; i < 10; i ++ {
		dataPointsFromDatabase[i].Time = now
		dataPointsFromDatabase[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}

	var counter = 0
	f := func(results mock.Arguments) () {
		for _, result := range results {
			switch v := result.(type) {
			case **int64:
				*v = &dataPointsFromDatabase[counter].Time

			case **float64:
				*v = &dataPointsFromDatabase[counter].Value.KWH
			}
		}
		counter ++
	}

	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(10)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)
	var aggregatedList []datatypes.Data
	lastValue := datatypes.Data{Time: -1,}
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, dataPointsFromDatabase[9], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t,dataPointsFromDatabase, aggregatedList, "Check if contains all values")

	kwh = 0
	 now = timestampSinceEpoch
	dataPointsFromDatabase = make([]datatypes.Data, 2)
	for i := 0; i < 2; i ++ {
		dataPointsFromDatabase[i].Time = now
		dataPointsFromDatabase[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}
	counter = 0
	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(2)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)
	aggregatedList = []datatypes.Data{}
	lastValue = datatypes.Data{Time: -1,}
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, dataPointsFromDatabase[1], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, dataPointsFromDatabase, aggregatedList, "Check if contains all values")


	kwh = 0
	now = timestampSinceEpoch
	dataPointsFromDatabase = make([]datatypes.Data, 1 )
	for i := 0; i < 1; i ++ {
		dataPointsFromDatabase[i].Time = now
		dataPointsFromDatabase[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}
	counter = 0
	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(1)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)
	aggregatedList = []datatypes.Data{}
	lastValue = datatypes.Data{Time: -1,}
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, dataPointsFromDatabase[0], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, append(dataPointsFromDatabase, dataPointsFromDatabase[0] ), aggregatedList, "Check if contains all values")

}

func TestWithNoPreciseAggregationPoints(t *testing.T) {

	var m, _ = &cassandra.IteratorMock{}, fmt.Errorf("e")

	var now int64 = timestampSinceEpoch +config.MILLI_SECONDS_BETWEEN_DATA_POINTS_LAYER_1
	var kwh float64
	dataPointsFromDatabase := make([]datatypes.Data, 10)
	for i := 0; i < 10; i ++ {
		dataPointsFromDatabase[i].Time = now
		dataPointsFromDatabase[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}

	var counter = 0
	f := func(results mock.Arguments) () {
		for _, result := range results {
			switch v := result.(type) {
			case **int64:
				*v = &dataPointsFromDatabase[counter].Time

			case **float64:
				*v = &dataPointsFromDatabase[counter].Value.KWH
			}
		}
		counter ++
	}

	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(10)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)
	var aggregatedList []datatypes.Data
	lastValue := datatypes.Data{Time: -1,}
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)

	now = timestampSinceEpoch +config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	kwh  = 5
	aggregatedResultList := make([]datatypes.Data, 9)
	for i := 0; i < 9; i ++ {
		aggregatedResultList[i].Time = now
		aggregatedResultList[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}

	assert.Equal(t, dataPointsFromDatabase[9], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, append(append([]datatypes.Data{dataPointsFromDatabase[0]}, aggregatedResultList...), dataPointsFromDatabase[9] ), aggregatedList, "Check if contains all values")

	now = timestampSinceEpoch +config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	kwh  = 5
	aggregatedResultList = make([]datatypes.Data, 1)
	for i := 0; i < 1; i ++ {
		aggregatedResultList[i].Time = now
		aggregatedResultList[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}

	counter = 0
	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(2)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)
	aggregatedList = []datatypes.Data{}
	lastValue = datatypes.Data{Time: -1,}
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, dataPointsFromDatabase[1], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, append(append([]datatypes.Data{dataPointsFromDatabase[0]}, aggregatedResultList...), dataPointsFromDatabase[1] ), aggregatedList, "Check if contains all values")

	dataPointsFromDatabase = make([]datatypes.Data, 1 )
	for i := 0; i < 1; i ++ {
		dataPointsFromDatabase[i].Time = now
		dataPointsFromDatabase[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}
	counter = 0
	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(1)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)
	aggregatedList = []datatypes.Data{}
	lastValue = datatypes.Data{Time: -1,}
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, dataPointsFromDatabase[0], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, append([]datatypes.Data{dataPointsFromDatabase[0]}, dataPointsFromDatabase...) , aggregatedList, "Check if contains all values")

}

func TestWithNoPreciseAggregationTimestampsWithPreviousTimestamp(t *testing.T) {

	var m, _ = &cassandra.IteratorMock{}, fmt.Errorf("e")

	var now int64 = timestampSinceEpoch +config.MILLI_SECONDS_BETWEEN_DATA_POINTS_LAYER_1
	var kwh float64 = 10
	dataPointsFromDatabase := make([]datatypes.Data, 10)
	for i := 0; i < 10; i ++ {
		dataPointsFromDatabase[i].Time = now
		dataPointsFromDatabase[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}

	var counter = 0
	f := func(results mock.Arguments) () {
		for _, result := range results {
			switch v := result.(type) {
			case **int64:
				*v = &dataPointsFromDatabase[counter].Time

			case **float64:
				*v = &dataPointsFromDatabase[counter].Value.KWH
			}
		}
		counter ++
	}

	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(10)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)

	var aggregatedList []datatypes.Data

	now = timestampSinceEpoch
	kwh  = 5
	aggregatedResultList := make([]datatypes.Data, 10)
	for i := 0; i < 10; i ++ {
		aggregatedResultList[i].Time = now
		aggregatedResultList[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}
	lastValue := datatypes.Data{Time : timestampSinceEpoch -config.MILLI_SECONDS_BETWEEN_DATA_POINTS_LAYER_1, }
	lastValue.Value.KWH = 0
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, dataPointsFromDatabase[9], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t, append(aggregatedResultList, dataPointsFromDatabase[9] ), aggregatedList, "Check if contains all values")


	counter = 0
	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(2)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)

	now = timestampSinceEpoch
	kwh  = 5
	aggregatedResultList = make([]datatypes.Data, 2)
	for i := 0; i < 2; i ++ {
		aggregatedResultList[i].Time = now
		aggregatedResultList[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}

	aggregatedList = []datatypes.Data{}
	lastValue = datatypes.Data{Time : timestampSinceEpoch -config.MILLI_SECONDS_BETWEEN_DATA_POINTS_LAYER_1, }
	lastValue.Value.KWH = 0
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, dataPointsFromDatabase[1], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t,  append(aggregatedResultList, dataPointsFromDatabase[1] ), aggregatedList, "Check if contains all values")


	counter = 0
	m.On("Scan", mock.Anything, mock.Anything).Run(f).Return(true).Times(1)
	m.On("Scan", mock.Anything, mock.Anything).Return(false).Times(1)
	aggregatedResultList = make([]datatypes.Data, 1 )
	now = timestampSinceEpoch
	kwh  = 5
	for i := 0; i < 1; i ++ {
		aggregatedResultList[i].Time = now
		aggregatedResultList[i].Value.KWH = kwh
		kwh += 10
		now += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
	}
	aggregatedList = []datatypes.Data{}
	lastValue = datatypes.Data{Time : timestampSinceEpoch -config.MILLI_SECONDS_BETWEEN_DATA_POINTS_LAYER_1, }
	lastValue.Value.KWH = 0
	lastValue, _ = DownsSampleRawDataToAggr1(&aggregatedList, lastValue, m)
	assert.Equal(t, dataPointsFromDatabase[0], lastValue, "Last value is datapoint with time 0")
	assert.Equal(t,   append(aggregatedResultList, dataPointsFromDatabase[0] ) , aggregatedList, "Check if contains all values")

}