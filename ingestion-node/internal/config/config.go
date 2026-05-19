package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

// holds the config/config.yaml data
type Config struct {
	LogFile string `yaml:"logFile"`
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

	return cfg
}

var cfg *Config

func init() {
	cfg = readConfigYaml()
}

func Get() *Config {
	return cfg
}
