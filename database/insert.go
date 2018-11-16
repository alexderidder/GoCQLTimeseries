package database

import (
	"../model"
	"../server"
	"github.com/gocql/gocql"
)

// i is not a good variable name! You know this!
func Insert(i *model.InsertJSON) model.Error {
	// ugly use of internals. What happens if the session disappears?
	batch := server.DbConn.Session.NewBatch(gocql.LoggedBatch)

	// why batch2? what does this do?
	batch2 := server.DbConn.Session.NewBatch(gocql.LoggedBatch)
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
			// a power factor without Watt is useless, the other way around could be plausible but we'll force these to be in pairs.
			// this case can be ignored.
			batch.Query("INSERT INTO w_and_pf_by_id_and_time (id, time, pf) VALUES (?, ?, ?)", i.StoneID.Value, data.Time, data.PowerFactor.Value)
		} else {
			// else?
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

		// whats the point of this? A comment would explain this a bit more.
		// The usage of something called batch2 is unclear
		if batch2.Size()%server.DbConn.BatchSize == 0 {
			err := server.DbConn.Session.ExecuteBatch(batch2)
			if err != nil {
				return model.Error{100, err.Error()}
			}
			batch2 = server.DbConn.Session.NewBatch(gocql.LoggedBatch)
		}
	}

	// why do we execute batch here again?
	err = server.DbConn.Session.ExecuteBatch(batch)
	if err != nil {
		return model.Error{100, err.Error()}
	}
	return model.NoError
}
