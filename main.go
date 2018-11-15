package main

import (
	"CrownstoneServer/server"
	"CrownstoneServer/server/tcp_server/connector"
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Server struct {
		IPAddress string `json:"ip-address"`
		Port      string `json:"port"`
		Certs     struct {
			Directory string `json:"directory"`
			Pem       string `json:"pem"`
			Key       string `json:"key"`
		} `json:"certs"`
		Messages struct {
			Timeout    uint32 `json:"timeout"`
			BufferSize uint32 `json:"buffer_size"`
		} `json:"messages"`
	} `json:"server"`
	Database struct {
		IPAddresses   []string `json:"ip-addresses"`
		Keyspace      string   `json:"keyspace"`
		ReconnectTime uint32   `json:"reconnect_time"`
		BatchSize     uint32   `json:"batch_size"`
	} `json:"database"`
}


func main() {
	config := readConfigFile()

	go connectCassandra(config)
	defer server.DbConn.Session.Close()
	startTLSServer(config)
}

func connectCassandra(config Configuration){
	ipAddresses := config.Database.IPAddresses
	keyspace := config.Database.Keyspace
	batchSize := int(config.Database.BatchSize)
	reconnectTime:= int(config.Database.ReconnectTime)
	server.ConnectCassandra(ipAddresses, keyspace, batchSize, reconnectTime)
}

func readConfigFile() Configuration {
	var config Configuration
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error parsing server config:", err)
		os.Exit(0)
	}
	return config
}

func startTLSServer(config Configuration){
	key := config.Server.Certs.Directory + config.Server.Certs.Pem
	cert := config.Server.Certs.Directory + config.Server.Certs.Key
	hostAndPort := config.Server.IPAddress + config.Server.Port
	timeOut := config.Server.Messages.Timeout
	bufferSize := config.Server.Messages.BufferSize
	connector.StartServerMode(key, cert, hostAndPort, timeOut, bufferSize)
}

