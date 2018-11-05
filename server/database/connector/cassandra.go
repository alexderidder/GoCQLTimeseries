package cassandra

import (
	"CrownstoneServer/server/config"
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

var Session *gocql.Session

var cluster *gocql.ClusterConfig

func StartCassandra() {
	clusterConfig()
	connect()
}

func connect() {
	Reconnect(true)

	fmt.Println("cassandra connection")
}
func Reconnect(firstTime bool) {
	var err error
	for {
		if !firstTime{
			time.Sleep(time.Duration(config.Config.Database.ReconnectTime) * time.Second)
		}


		Session, err = cluster.CreateSession()
		if err != nil {
			//TODO: Printing 3 lines with the same error when keyspace doesn't exists
			firstTime = false
			fmt.Println(err)
		} else {
			return
		}
	}

}
func clusterConfig() {
	var ipAddresses []string
	for _, value := range config.Config.Database.Clusters {
		ipAddresses = append(ipAddresses, value.IPAddress)
	}
	cluster = gocql.NewCluster(ipAddresses...)
	cluster.Keyspace = config.Config.Database.Keyspace
}
