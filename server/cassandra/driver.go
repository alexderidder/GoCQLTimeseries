package cassandra

import (
	"../../model"
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

type db struct {
	Session       *gocql.Session
	BatchSize     int
	ReconnectTime int
	cluster       *gocql.ClusterConfig
}


var dbConn *db

func ConnectCassandra(ipAddresses []string, keyspace string, batchSize int, reconnectTime int) {
	dbConn = &db{}
	dbConn.cluster = gocql.NewCluster(ipAddresses...)
	dbConn.cluster.Keyspace = keyspace
	dbConn.BatchSize = batchSize
	dbConn.ReconnectTime = reconnectTime
	connect()
	fmt.Println("Database connected")
}

func Close() {
	dbConn.Session.Close()
}

func connect() {
	var err error
	for {
		dbConn.Session, err = dbConn.cluster.CreateSession()
		if err != nil {
			//TODO: Printing 3 lines with the same error when keyspace doesn't exists
			fmt.Println(err)
			time.Sleep(time.Duration(dbConn.ReconnectTime) * time.Second)
		} else {
			fmt.Println("Database is connected")
			return
		}
	}
}

func reconnect(tries uint32) model.Error {
	var err error
	for i := uint32(0); i < tries; i++ {
		dbConn.Session, err = dbConn.cluster.CreateSession()
		if err != nil {
			//TODO: Printing 3 lines with the same error when keyspace doesn't exists
			fmt.Println(err)
			time.Sleep(time.Duration(dbConn.ReconnectTime) * time.Second)
		} else {
			fmt.Println("Database is connected")
			return model.NoError
		}
	}
	return model.ServerNoCassandra
}
//@Alex What is a good number of tries before error return ?
func checkConnection() model.Error {
	if dbConn.Session.Closed() {
		reconnect(3)
	}
	return model.NoError
}

func Query(stmt string, values ...interface{}) (*gocql.Iter, model.Error) {
	if err := checkConnection(); !err.IsNull() {
		return nil, err

	}
	return dbConn.Session.Query(stmt, values...).Iter(), model.NoError
}

func CreateBatch() (*gocql.Batch, model.Error) {
	if err := checkConnection(); !err.IsNull() {
		return nil, err

	}
	return dbConn.Session.NewBatch(gocql.LoggedBatch), model.NoError
}

func ExecuteBatch(batch *gocql.Batch) model.Error {
	if err := checkConnection(); !err.IsNull() {
		return err
	}
	if batch.Size() > 0 {
		err := dbConn.Session.ExecuteBatch(batch)
		if err != nil {
			return model.Error{100, err.Error()}
		}
	}
	return model.NoError
}

func AddQueryToBatch(batch *gocql.Batch, stmt string, values ...interface{}) model.Error {
	if err := checkConnection(); !err.IsNull() {
		return err
	}
	batch.Query(stmt, values...)
	if batch.Size()%dbConn.BatchSize == 0 {
		err := dbConn.Session.ExecuteBatch(batch)
		if err != nil {
			return model.Error{100, err.Error()}
		}
		*batch = *dbConn.Session.NewBatch(gocql.LoggedBatch)

	}
	return model.NoError

}
