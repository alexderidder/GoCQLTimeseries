package server

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

type DB struct {
	Session *gocql.Session
	BatchSize int
	ReconnectTime int
	cluster *gocql.ClusterConfig
}

var DbConn = DB{}

func ConnectCassandra(ipAddresses []string,  keyspace string, batchSize int, reconnectTime int){
	DbConn.cluster = gocql.NewCluster(ipAddresses...)
	DbConn.cluster.Keyspace = keyspace
	DbConn.BatchSize = batchSize
	DbConn.ReconnectTime = reconnectTime

	Reconnect(true)
	fmt.Println("Database connected")
}

func Reconnect(firstTime bool) {
	var err error
	for {
		if !firstTime{
			time.Sleep(time.Duration(DbConn.ReconnectTime) * time.Second)
		}
		DbConn.Session, err = DbConn.cluster.CreateSession()
		if err != nil {
			//TODO: Printing 3 lines with the same error when keyspace doesn't exists
			firstTime = false
			fmt.Println(err)
		} else {
			return
		}
	}

}
