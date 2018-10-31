package config

import (
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
			Timeout    int `json:"timeout"`
			BufferSize uint32 `json:"buffer_size"`
		} `json:"messages"`
	} `json:"server"`
	Database struct {
		Clusters []struct {
			Name      string `json:"name"`
			IPAddress string `json:"ip-address"`
		} `json:"clusters"`
		Keyspace      string `json:"keyspace"`
		ReconnectTime uint32    `json:"reconnect_time"`
	} `json:"database"`
}

var Config Configuration

func init() {
	file, _ := os.Open("server/config/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error parsing server config:", err)
	}

}
