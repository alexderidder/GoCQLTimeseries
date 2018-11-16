package server

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

type DB struct {
	Session       *gocql.Session
	BatchSize     int
	ReconnectTime int
	cluster       *gocql.ClusterConfig
}

var DbConn = DB{}

func ConnectCassandra(ipAddresses []string, keyspace string, batchSize int, reconnectTime int) {
	DbConn.cluster = gocql.NewCluster(ipAddresses...)
	DbConn.cluster.Keyspace = keyspace
	DbConn.BatchSize = batchSize
	DbConn.ReconnectTime = reconnectTime

	// ugly to call reconnect for initial connection
	Reconnect(true)
	fmt.Println("Database connected")
}

// Rename to Connect, possibly have a reconnect method that will first sleep, then call connect. EDIT Ah. I see now why its like this..
func Reconnect(firstTime bool) {
	var err error
	// why is there a for here? Is this a retry loop? EDIT Ah. I see now why its like this.. What happens if there is already a connection?
	for {
		if !firstTime { // if reconnect and connect are separated, this will not be in the connect, making that more readable. EDIT: see below
			time.Sleep(time.Duration(DbConn.ReconnectTime) * time.Second)
		}
		DbConn.Session, err = DbConn.cluster.CreateSession()
		if err != nil {
			//TODO: Printing 3 lines with the same error when keyspace doesn't exists
			firstTime = false
			fmt.Println(err)
		} else {
			return // break seems more appropriate to get out of a for loop, you might do some bookkeeping below this reconnecting loop once the connection is established.
		}
	}

}

// I'd probably perfer this kind of setup, or maybe recursive. Downside with recursion is that you will hit a stack limit somewhere.
//func Connect() {
//	var err error
//	for {
//		DbConn.Session, err = DbConn.cluster.CreateSession()
//      // Retry the connection
//		if err != nil {
//			fmt.Println(err)
//			time.Sleep(time.Duration(DbConn.ReconnectTime) * time.Second)
//		} else {
//			break
//		}
//	}
//}
