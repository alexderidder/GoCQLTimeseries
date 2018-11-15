package database

import (
	"../model"
	"../server"
	"github.com/gocql/gocql"
)

func Insert(i *model.InsertJSON) model.Error {

	batch := server.DbConn.Session.NewBatch(gocql.LoggedBatch)
	batch2 :=  server.DbConn.Session.NewBatch(gocql.LoggedBatch)
	var err error
	//TODO: index isnt constant anymore
	for _, data := range i.Data {

		if data.Watt.Valid {
			if data.PowerFactor.Valid {
				batch.Query("INSERT INTO w_and_pf_by_id_and_time (id, time, w, pf) VALUES (?, ?, ?, ?)", i.StoneID.Value, data.Time, data.Watt.Value, data.PowerFactor.Value)
			} else {
				batch.Query("INSERT INTO w_and_pf_by_id_and_time (id, time, w) VALUES (?, ?, ?)", i.StoneID.Value, data.Time, data.Watt.Value)
			}
		} else if data.PowerFactor.Valid {
			batch.Query("INSERT INTO w_and_pf_by_id_and_time (id, time, pf) VALUES (?, ?, ?)", i.StoneID.Value, data.Time, data.PowerFactor.Value)
		} else {

		}

		if batch.Size()%server.DbConn.BatchSize == 0 {
			err := server.DbConn.Session.ExecuteBatch(batch)
			if err != nil {
				return model.Error{100, err.Error()}
			}
			batch = server.DbConn.Session.NewBatch(gocql.LoggedBatch)
		}

		if data.KWH.Valid {
			batch2.Query("INSERT INTO kwh_by_id_and_time (id, time, kwh) VALUES (?, ?, ?)", i.StoneID.Value, data.Time, data.KWH.Value)
		} else {
			continue
		}

		if batch2.Size()%server.DbConn.BatchSize == 0 {
			err := server.DbConn.Session.ExecuteBatch(batch2)
			if err != nil {
				return model.Error{100, err.Error()}
			}
			batch2 = server.DbConn.Session.NewBatch(gocql.LoggedBatch)
		}
	}
	err = server.DbConn.Session.ExecuteBatch(batch)
	if err != nil {
		return model.Error{100, err.Error()}
	}
	return model.NoError
}
