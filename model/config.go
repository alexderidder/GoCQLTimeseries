package model

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Server struct {
		IPAddress string `json:"ip-address"`
		Port      string `json:"port"`
		Certs     struct {
			Directory string `json:"directory"`
			Cert       string `json:"cert"`
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

func DecodeConfigJSON(path string) (Configuration, error){
	var config Configuration
	file, _ := os.Open(path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err := decoder.Decode(&config)
	return config, err

}
