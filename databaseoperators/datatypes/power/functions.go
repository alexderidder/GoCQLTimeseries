package power

import (
	"GoCQLTimeSeries/config"
	"GoCQLTimeSeries/datatypes"
	"GoCQLTimeSeries/server/cassandra"
	"fmt"
	"github.com/gocql/gocql"
)

type AggregatePower struct {
	sumPF  float64
	sumW  float64
	timestamp int64
	count float64
}

func DownsSampleRawDataToAggr1(dataList *[]datatypes.Data, test AggregatePower, iterator cassandra.Iterator) (AggregatePower, datatypes.Error) {
	var currentSum AggregatePower
	var timestampForAggregatedDataPoint int64

	var timestamp *int64
	var wattage *float32
	var powerfactor *float32
	if test.timestamp != -1 {
		currentSum = test
		timestampForAggregatedDataPoint = test.timestamp - ( test.timestamp % config.MILLISECONDS_INTERVAL_FOR_POWER_AGGREGATION) + config.LAYER_1_DATA_POINTS_INTERVAL

	} else {
		for iterator.Scan(&timestamp, &wattage, &powerfactor) {
				if timestamp != nil {
					if *timestamp < timestampForAggregatedDataPoint {
						currentSum.sumW += float64(*wattage)
						currentSum.sumPF += float64(*powerfactor)
						currentSum.count++
						currentSum.timestamp = *timestamp
					} else {
						if currentSum.count != 0 {
							data := datatypes.Data{Time: timestampForAggregatedDataPoint - config.MILLISECONDS_INTERVAL_FOR_POWER_AGGREGATION}
							data.Value.Wattage = float32(currentSum.sumW / (currentSum.count))
							data.Value.PowerFactor = float32(currentSum.sumPF / (currentSum.count))
							*dataList = append(*dataList, data)

							currentSum.sumW = float64(*wattage)
							currentSum.sumPF = float64(*powerfactor)
							currentSum.count = 1
							timestampForAggregatedDataPoint =  test.timestamp - ( test.timestamp % config.MILLISECONDS_INTERVAL_FOR_POWER_AGGREGATION) + config.LAYER_1_DATA_POINTS_INTERVAL
						}
					}
				}
			}

	}

	if err := iterator.Close(); err != nil {
		error := datatypes.CassandraIterator
		error.Message = err.Error()
		return currentSum, error

	}
	return currentSum, datatypes.NoError

}

func GetAlreadyDownSampledData(iterator *gocql.Iter, dataList *[]datatypes.Data) {
	var timeOfRow *int64
	var w, pf *float32
	for iterator.Scan(&timeOfRow, &w, &pf) {
		if timeOfRow != nil {
			var data = datatypes.Data{Time: *timeOfRow}
			if w != nil && pf != nil {
				data.Value.Wattage = *w
				data.Value.PowerFactor = *pf
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
				cassandra.AddQueryToBatchAndExecuteBatchTillSuccess(batch, "INSERT INTO w_and_pf_by_id_and_time_layer1  (id, time_bucket, time, w, pf) VALUES (?,?, ?, ?, ?)", dataWithStoneID.StoneID, dataWithBucket.Bucket, dataPoint.Time, dataPoint.Value.Wattage, dataPoint.Value.PowerFactor)
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
