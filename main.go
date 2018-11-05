package main

import (
	"CrownstoneServer/server/database/connector"
	"CrownstoneServer/server/tcp_server/connector"
)

func main() {
	go cassandra.StartCassandra()
	defer cassandra.Session.Close()
	connector.StartServerMode()
}
