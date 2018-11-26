package main

import (
	"./model"
	"./server/cassandra"
	"./server/socket/tls"
)

//Reads the config file, then gocql driver connects to Cassandra with config values in a new Thread/GoRoutine. The main blocks with the StartTLS server, it accepts a connection and adds a Thread to communicate with this connection.

func main() {

	config, err := model.DecodeConfigJSON("conf.json")
	if err != nil {
		//Close program because config attributes are needed
		panic(err)
	}

	go cassandra.ConnectCassandra(config.Database.IPAddresses, config.Database.Keyspace, int(config.Database.BatchSize), int(config.Database.ReconnectTime))

	//Blocking
	tls.StartTLSServer(config.Server.Certs.Directory+config.Server.Certs.Pem, config.Server.Certs.Directory+config.Server.Certs.Key, config.Server.IPAddress+config.Server.Port, config.Server.Messages.Timeout, config.Server.Messages.BufferSize)
}
