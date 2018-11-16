package main

import (
	"./server"
	"./server/tcp_server/connector"
	"encoding/json"
	"fmt"
	"os"
)

// Since MAIN is not describing function name, add a bit of comment on what this function does
// I understand main is a mandatory name, but you can describe what the intended use of this method is.

// why are we providing the --mode server?
func main() {

	// I'd probably make a function that will just get the config or error
	// that way, when reading the main loop, we're not going through 10 lines of arbitrary code to parse a configuration
	// but we can immediately see what is happening. This is nice for other people reading the code.
	var config Configuration
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error parsing server config:", err)
		os.Exit(0) // why not panic?
	}

	// where does this goroutine come back to the main thread? The server close is deferred but it does not perse wait for this goroutine.
	go server.ConnectCassandra(config.Database.IPAddresses, config.Database.Keyspace, int(config.Database.BatchSize), int(config.Database.ReconnectTime))

	// I don't like the going into the internals. You expect the session to be here when the defer is fired. Is not implied.
	// Try to design all layers as if they have an API to the other layers using them.
	// the main should not need to know the internals of the server.
	defer server.DbConn.Session.Close()

	// the success of this main function depends on that this method is blocking. That is the only reason the deferred close does not have to wait on the connect. This sort of behaviour will lead to big crashes while refactoring.
	// If your server was decoupled (so all server logic would happen in server), the server could guard for this.
	connector.StartServerMode(config.Server.Certs.Directory+config.Server.Certs.Pem, config.Server.Certs.Directory+config.Server.Certs.Key, config.Server.IPAddress+config.Server.Port, config.Server.Messages.Timeout, config.Server.Messages.BufferSize)
}

// I'd suggest a folder containing types so that this will not have to be here.
// Typedefs can be combined a bit to make the actual function files a bit smaller
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
