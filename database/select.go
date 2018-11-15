package database

import (
	"../model"
	"../server"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"time"
)

func Select(s *model.RequestSelectJSON ) (model.ResponseSelectJSON, model.Error) {
	response := model.ResponseSelectJSON{s.StartTime, s.EndTime, s.Interval, []model.Stone{}}

	var timeValues []interface{}
	var timeQuery string
	if !s.StartTime.IsZero() && !s.EndTime.IsZero() {
		timeQuery = ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, s.StartTime)
		timeValues = append(timeValues, s.EndTime)
	}
	var queryValues []interface{}

	for _, stoneID := range s.StoneIDs {
		queryValues = append([]interface{}{}, stoneID.Value)
		queryValues = append(queryValues, timeValues...)

		stone := model.Stone{}
		stone.StoneID = stoneID.Value
		stone.Fields = []model.Field{}
		var iterator *gocql.Iter

		switch s.Types[0] {
		case model.UnitWAndpf:
			iterator = server.DbConn.Session.Query("SELECT time, " + model.UnitW + ", " + model.Unitpf + " FROM w_and_pf_by_id_and_time WHERE id = ?" + timeQuery ,  queryValues...).Iter()
			var timeOfRow time.Time
			var w, pf *float32
			var wList, pfList []model.Data
			for iterator.Scan(&timeOfRow, &w, &pf) {
				if w != nil {
					wList = append(wList, model.Data{timeOfRow, *w})
				}
				if pf != nil {
					pfList = append(pfList, model.Data{timeOfRow, *pf})
				}
			}
			if err := iterator.Close(); err != nil {
				log.Fatal(err)
				//TODO: research if error code is needed
			}
			stone.Fields = append(stone.Fields, model.Field{ model.UnitW, wList})
			stone.Fields = append(stone.Fields,  model.Field{ model.Unitpf, pfList})
		case  model.UnitW:
			iterator = server.DbConn.Session.Query("SELECT time, " + model.UnitW + " FROM w_and_pf_by_id_and_time WHERE id = ?" + timeQuery,  queryValues...).Iter()
			measurements, err := iterateStreamWithOneFloat32PerRow(iterator, s.Interval)
			if err != nil {
				return response, model.Error{300, err.Error()}
			}
			stone.Fields = append(stone.Fields, model.Field{model.UnitW, measurements})
		case  model.Unitpf:
			iterator = server.DbConn.Session.Query("SELECT time, " + model.Unitpf + " FROM w_and_pf_by_id_and_time WHERE id = ?" + timeQuery,  queryValues...).Iter()
			measurements, err := iterateStreamWithOneFloat32PerRow(iterator, s.Interval)
			if err != nil {
				return response, model.Error{300, err.Error()}
			}

			stone.Fields = append(stone.Fields, model.Field{model.Unitpf, measurements})
		}

		switch s.Types[1] {
		case  model.UnitkWh:
			iterator = server.DbConn.Session.Query("SELECT time, " + model.UnitkWh + " FROM kwh_by_id_and_time WHERE id = ?" + timeQuery, queryValues...).Iter()
			measurements, err := iterateStreamWithOneFloat32PerRow(iterator, s.Interval)
			if err != nil {
				return response,  model.Error{300, err.Error()}
			}
			stone.Fields = append(stone.Fields, model.Field{model.UnitkWh, measurements})
		}
		response.Stones = append(response.Stones, stone)
	}
	return response, model.NoError

}

func  iterateStreamWithOneFloat32PerRow(iterator *gocql.Iter, interval uint32) ([]model.Data, error) {
	var value *float32
	var measurementList []model.Data
	var timeOfRow time.Time

	for iterator.Scan(&timeOfRow, &value) {
		if value != nil{
			measurementList = append(measurementList, model.Data{timeOfRow, *value})
		}
	}

	if err := iterator.Close(); err != nil {
		fmt.Println(err)
		return measurementList, err
		//TODO: error code is needed
	}

	return measurementList, nil
}

