package database

import (
	"CrownstoneServer/model"
	"CrownstoneServer/server"
	"github.com/gocql/gocql"
	"log"
	"time"
)

func Delete(d *model.DeleteJSON) (model.Error) {

	var queryTimePart string
	var timeValues []interface{}
	timeValues = append(timeValues, d.StoneID.Value)
	if !d.StartTime.IsZero() && !d.EndTime.IsZero() {
		queryTimePart = ` AND time >= ? AND time <= ? `
		timeValues = append(timeValues, d.StartTime)
		timeValues = append(timeValues, d.EndTime)
	}

	var values string
	var error model.Error
	switch d.Types[0] {
	case model.UnitWAndpf:
		values = model.UnitW + " = null, " + model.Unitpf + " = null"
	case model.UnitW:
		values = model.UnitW + " = null"
	case model.Unitpf:
		values = model.Unitpf + " = null"
	default:
		goto Skip
	}
	values = "UPDATE w_and_pw_by_id_and_time SET " + values + " WHERE id = ? AND time = ?"

	error = selectAndInsert("SELECT time FROM w_and_pw_by_id_and_time WHERE id = ?"+queryTimePart, values, timeValues)
	if !error.IsNull() {
		return error
	}
Skip:
	switch d.Types[1] {
	case model.UnitkWh:
		error = selectAndInsert("SELECT time FROM kwh_by_id_and_time WHERE id = ?"+queryTimePart, "UPDATE kwh_by_id_and_time SET "+model.UnitkWh+" = null WHERE id = ? AND time = ?", timeValues)
		if !error.IsNull() {
			return error
		}
	}
	return model.NoError
}

func selectAndInsert(selectQuery string, insertQuery string, values []interface{}) model.Error {
	var err error
	iterator := server.DbConn.Session.Query(selectQuery, values...).Iter()
	var timeOfRow time.Time
	var timeArray []time.Time
	for iterator.Scan(&timeOfRow) {
		timeArray = append(timeArray, timeOfRow)
	}

	if err := iterator.Close(); err != nil {
		log.Fatal(err)
		//TODO: research if error code is needed
	}
	batch := server.DbConn.Session.NewBatch(gocql.LoggedBatch)
	for index, valueTime := range timeArray {
		batch.Query(insertQuery, values[0], valueTime)

		if index%server.DbConn.BatchSize == 0 {
			err := server.DbConn.Session.ExecuteBatch(batch)
			if err != nil {

				return model.Error{100, err.Error()}
			}
			batch = server.DbConn.Session.NewBatch(gocql.LoggedBatch)
		}
	}

	if batch.Size() > 0 {
		err = server.DbConn.Session.ExecuteBatch(batch)
		if err != nil {
			return model.Error{100, err.Error()}
		}
	}
	return model.NoError
}
