package cassandra

import (
	"GoCQLTimeSeries/datatypes"
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
	//Blocking, check every second if session is closed. Then connect
	for {
		if dbConn.Session.Closed(){
			connect()
			fmt.Println("Connection closed with cassandra, trying to reconnect")
		}		else{
			time.Sleep(time.Second)
		}
	}
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
			time.Sleep(time.Duration(dbConn.ReconnectTime))
		} else {
			fmt.Println("Database is connected")
			return
		}
	}
}



func checkConnection() datatypes.Error {
	if dbConn.Session.Closed() {
		return datatypes.ServerNoCassandra
	}
	return datatypes.NoError
}

func Query(stmt string, values ...interface{}) (*gocql.Iter, datatypes.Error) {
	if err := checkConnection(); !err.IsNull() {
		return nil, err
	}
	return dbConn.Session.Query(stmt, values...).Iter(), datatypes.NoError
}

func ExecQuery(stmt string, values ...interface{}) ( datatypes.Error) {
	if err := checkConnection(); !err.IsNull() {
		return  err

	}
	fmt.Println(stmt, values)
	err := dbConn.Session.Query(stmt, values...).Exec()
	if err != nil {
		return datatypes.ExecuteCassandra
	}
	return datatypes.NoError
}

func CreateBatch() (*gocql.Batch, datatypes.Error) {
	if err := checkConnection(); !err.IsNull() {
		return nil, err

	}
	return dbConn.Session.NewBatch(gocql.LoggedBatch), datatypes.NoError
}

func ExecuteBatch(batch *gocql.Batch) datatypes.Error {
	if err := checkConnection(); !err.IsNull() {
		return err
	}
	if batch.Size() > 0 {
		err := dbConn.Session.ExecuteBatch(batch)
		if err != nil {
			return datatypes.Error{100, err.Error()}
		}
		batch = nil
	}
	return datatypes.NoError
}

func AddQueryToBatchAndExecuteWhenBatchMax(batch *gocql.Batch, stmt string, values ...interface{}) datatypes.Error {
	if err := checkConnection(); !err.IsNull() {
		return err
	}
	batch.Query(stmt, values...)
	if batch.Size()%dbConn.BatchSize == 0 {
		return ExecuteAndClearBatch(batch)
	}
	return datatypes.NoError
}

func ExecuteAndClearBatch(batch *gocql.Batch) datatypes.Error {
	if batch.Size() > 0 {

		err := dbConn.Session.ExecuteBatch(batch)
		if err != nil {
			return datatypes.Error{100, err.Error()}
		}
		*batch = *dbConn.Session.NewBatch(gocql.LoggedBatch)
	}
	return datatypes.NoError
}


func ExecuteBatchWithoutError(batch *gocql.Batch) {
	if err := checkConnection(); !err.IsNull() {
		time.Sleep(time.Duration(dbConn.ReconnectTime))

	}
	for {
		if err := dbConn.Session.ExecuteBatch(batch); err != nil {
			time.Sleep(time.Duration(dbConn.ReconnectTime))
			continue
		}
		break;
	
	}
}

func AddQueryToBatchAndExecuteBatchTillSuccess(batch *gocql.Batch, stmt string, values ...interface{})  {
	for  err := checkConnection(); !err.IsNull(); {
		time.Sleep(time.Duration(dbConn.ReconnectTime))
	}
	batch.Query(stmt, values...)
	if batch.Size()%dbConn.BatchSize == 0 {
		go ExecuteAndClearBatch(batch)
	}

}
