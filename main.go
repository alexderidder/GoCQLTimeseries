package main

import (
	"GoCQLSockets/server/database/connector"
	"GoCQLSockets/server/tcp_server/connector"
)

func main() {
	go cassandra.StartCassandra()
	CassandraSession := cassandra.Session
	defer CassandraSession.Close()
	connector.StartServerMode()
}
