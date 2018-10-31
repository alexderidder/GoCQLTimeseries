package main

import (
	"CrownstoneServer/server/database/connector"
	"CrownstoneServer/server/tcp_server/connector"
)

func main() {
	go cassandra.StartCassandra()
	CassandraSession := cassandra.Session
	defer CassandraSession.Close()
	connector.StartServerMode()
}
