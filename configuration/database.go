package configuration

import (
	"encoding/json"

	"io/ioutil"

	"github.com/JesusIslam/salestock/logger"
)

type DatabaseConfig struct {
	URL string `json:"url"`
}

func Database() *DatabaseConfig {
	data, err := ioutil.ReadFile(DefaultDirectory + "/database.json")
	if err != nil {
		logger.Fatal("Failed to load database configuration file: ", err)
	}

	databaseConfig := &DatabaseConfig{}
	err = json.Unmarshal(data, databaseConfig)
	if err != nil {
		logger.Fatal("Failed to unmarshal database configuration file: ", err)
	}

	return databaseConfig
}
