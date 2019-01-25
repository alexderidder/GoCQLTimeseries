package energy

import (
	"GoCQLTimeSeries/config"
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/server/cassandra"
	"GoCQLTimeSeries/util"
	"fmt"
	"github.com/gocql/gocql"
)

func DownsSampleRawDataToAggr1(dataList *[]datatypes.Data, lastDataPoint datatypes.Data, iterator cassandra.Iterator) (datatypes.Data, datatypes.Error) {

	var y1LastValue *float64
	var x1LastTimestamp *int64
	var y0PreviousValue = lastDataPoint.Value.KWH
	var x0PreviousTimestamp = lastDataPoint.Time
	var aggregatedXOfDataPoint int64

	if lastDataPoint.Time != -1 {
		y0PreviousValue = lastDataPoint.Value.KWH
		x0PreviousTimestamp = lastDataPoint.Time
		aggregatedXOfDataPoint = x0PreviousTimestamp - (x0PreviousTimestamp % config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION) + config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION

	} else {
		if iterator.Scan(&x1LastTimestamp, &y1LastValue) {
			y0PreviousValue = *y1LastValue
			x0PreviousTimestamp = *x1LastTimestamp
			aggregatedXOfDataPoint = x0PreviousTimestamp - (x0PreviousTimestamp % config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION) + config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
			data := datatypes.Data{Time: x0PreviousTimestamp}
			data.Value.KWH = y0PreviousValue
			*dataList = append(*dataList, data)
		} else {
			goto SkipIterator
		}

	}

	for iterator.Scan(&x1LastTimestamp, &y1LastValue) {

		if x1LastTimestamp != nil {

			// Initialize aggregatedXOfDataPoint to rounded value of x0PreviousTimestamp at 5 minutes
			if *x1LastTimestamp > aggregatedXOfDataPoint {
				// If new timestamp is not in 5 minute range of the aggregated x point store previous dataPoint
				if *x1LastTimestamp >= aggregatedXOfDataPoint+config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION {
					data := datatypes.Data{Time: x0PreviousTimestamp}
					data.Value.KWH = y0PreviousValue;
					*dataList = append(*dataList, data)

					data = datatypes.Data{Time: *x1LastTimestamp}
					data.Value.KWH = *y1LastValue;
					*dataList = append(*dataList, data)
					aggregatedXOfDataPoint = *x1LastTimestamp - (*x1LastTimestamp % config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION) + config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
				} else {
					//Linear interpolation
					kwh, _ := util.LinearInterpolation(x0PreviousTimestamp, *x1LastTimestamp, y0PreviousValue, *y1LastValue, aggregatedXOfDataPoint)
					data := datatypes.Data{Time: aggregatedXOfDataPoint}
					data.Value.KWH = kwh
					*dataList = append(*dataList, data)
					aggregatedXOfDataPoint += config.MILLI_SECONDS_INTERVAL_FOR_ENERGY_AGGREGATION
				}

			}

		}
		y0PreviousValue = *y1LastValue
		x0PreviousTimestamp = *x1LastTimestamp

	}

SkipIterator:

	if err := iterator.Close(); err != nil {
		error := datatypes.CassandraIterator
		error.Message = err.Error()
		return datatypes.Data{}, error

	}

	data := datatypes.Data{Time: x0PreviousTimestamp}
	data.Value.KWH = y0PreviousValue;
	if x0PreviousTimestamp != -1 {
		*dataList = append(*dataList, data)
	}
	return data, datatypes.NoError
}

func GetAlreadyDownSampledData(iterator *gocql.Iter, dataList *[]datatypes.Data) {
	var kWh *float64
	var timeOfRow *int64
	for iterator.Scan(&timeOfRow, &kWh) {
		if timeOfRow != nil {
			data := datatypes.Data{Time: *timeOfRow}
			if *kWh != 0 {
				data.Value.KWH = *kWh
				*dataList = append(*dataList, data)
			}
		}
	}
}

func InsertAggregatedDataInCassandra(dataPerStoneIDs []datatypes.StoneIDsWithBucketsWithDataPoints) {
	batch, error := cassandra.CreateBatch()
	if !error.IsNull() {
		fmt.Println(error)
		return
	}

	for _, dataWithStoneID := range dataPerStoneIDs {
		for _, dataWithBucket := range dataWithStoneID.BucketsWithDataPoints {

			for _, dataPoint := range dataWithBucket.EnergyDataPoints {
				cassandra.AddQueryToBatchAndExecuteBatchTillSuccess(batch, "INSERT INTO kWh_by_id_and_time_in_layer2 (id, time_bucket, time, kwh) VALUES (?,?, ?, ?)", dataWithStoneID.StoneID, dataWithBucket.Bucket, dataPoint.Time, dataPoint.Value.KWH)
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
