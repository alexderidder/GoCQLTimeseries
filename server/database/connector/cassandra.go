package cassandra

import (
	"GoCQLSockets/server/config"
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
	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		fmt.Println(err)
		Reconnect()
	}
	fmt.Println("cassandra connection")
}
func Reconnect() {
	var err error
	for {
		time.Sleep(time.Duration(config.Config.Database.ReconnectTime) * time.Second)
		Session, err = cluster.CreateSession()
		if err != nil {
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
