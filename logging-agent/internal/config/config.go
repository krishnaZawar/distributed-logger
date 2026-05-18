package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

// holds the config/config.yaml data
type Config struct {
	LogFiles         []string           `yaml:"logFiles"`
	AgentLogFile     string             `yaml:"agentLogFile"`
	DeliveryDetails  LogDeliveryDetails `yaml:"logDelivery"`
	LogReadBatchSize int                `yaml:"logReadBatchSize"`
}

type LogDeliveryDetails struct {
	Method             string `yaml:"method"`
	Endpoint           string `yaml:"endpoint"`
	ExpectedStatusCode int    `yaml:"expectedStatusCode"`
}

func readConfigYaml() *Config {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	cfg := &Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}

	if cfg.DeliveryDetails.Method == "" || cfg.DeliveryDetails.Endpoint == "" {
		panic("Error: log delivery details should be configured")
	}

	return cfg
}

var cfg *Config

func init() {
	cfg = readConfigYaml()
}

func Get() *Config {
	return cfg
}
