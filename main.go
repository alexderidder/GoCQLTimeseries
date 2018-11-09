package main

import (
	"CrownstoneServer/server"
	"CrownstoneServer/server/tcp_server/connector"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
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

	go server.ConnectCassandra(config.Database.IPAddresses, config.Database.Keyspace, int(config.Database.BatchSize), int(config.Database.ReconnectTime) )

	defer server.DbConn.Session.Close()
	connector.StartServerMode(config.Server.Certs.Directory + config.Server.Certs.Pem, config.Server.Certs.Directory + config.Server.Certs.Key, config.Server.IPAddress + config.Server.Port, config.Server.Messages.Timeout,  config.Server.Messages.BufferSize)
}

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
		IPAddresses []string `json:"ip-addresses"`
		Keyspace      string `json:"keyspace"`
		ReconnectTime uint32    `json:"reconnect_time"`
		BatchSize uint32 `json:"batch_size"`
	} `json:"database"`
}