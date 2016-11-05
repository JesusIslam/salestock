package configuration

import (
	"encoding/json"

	"io/ioutil"

	"github.com/JesusIslam/salestock/logger"
)

type ServerConfig struct {
	Host string `json:"host"`
}

func Server() *ServerConfig {
	data, err := ioutil.ReadFile(DefaultDirectory + "/server.json")
	if err != nil {
		logger.Fatal("Failed to load server configuration file: ", err)
	}

	serverConfig := &ServerConfig{}
	err = json.Unmarshal(data, serverConfig)
	if err != nil {
		logger.Fatal("Failed to unmarshal server configuration file: ", err)
	}

	return serverConfig
}
