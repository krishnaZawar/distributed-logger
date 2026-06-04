package config

import (
	"os"

	"github.com/krishnaZawar/distributed-logger/logging-agent/internal/base"
	"go.yaml.in/yaml/v3"
)

// holds the configuration details configured in config.yaml
type Config struct {
	LogFiles     []string `yaml:"logFiles"`
	AgentLogFile string   `yaml:"agentLogFile"`

	LogDeliveryDetails LogDeliveryConfigurations `yaml:"logDelivery"`

	OffsetPersistenceDetails OffsetPersistenceConfigurations `yaml:"offsetPersistence"`

	LogCollectionDetails LogCollectionConfigurations `yaml:"logCollection"`
}

// holds log collection configuration details
type LogCollectionConfigurations struct {
	CollectionBatchSize  int `yaml:"batchSize"`
	CollectionIntervalMs int `yaml:"intervalMs"`
}

// holds the offset persistence configuration details
type OffsetPersistenceConfigurations struct {
	UpdateIntervalMs int                 `yaml:"updateIntervalMs"`
	RetryDetails     RetryConfigurations `yaml:"retryConfiguration"`
}

// holds the retry configuration details
type RetryConfigurations struct {
	TimeoutMs int `yaml:"timeoutMs"`
	Count     int `yaml:"count"`
}

// holds the log delivery configuration details
type LogDeliveryConfigurations struct {
	Method             string              `yaml:"method"`
	Endpoint           string              `yaml:"endpoint"`
	ExpectedStatusCode int                 `yaml:"expectedStatusCode"`
	RetryDetails       RetryConfigurations `yaml:"retryConfiguration"`
}

// creates a new config with the default values
func newConfig() *Config {
	return &Config{
		LogFiles:     []string{},
		AgentLogFile: "",

		LogDeliveryDetails: LogDeliveryConfigurations{
			Method:             "",
			Endpoint:           "",
			ExpectedStatusCode: 0,
			RetryDetails: RetryConfigurations{
				TimeoutMs: base.DefaultDeliveryRetryTimeoutMs,
				Count:     base.DefaultDeliveryRetryCount,
			},
		},

		OffsetPersistenceDetails: OffsetPersistenceConfigurations{
			UpdateIntervalMs: base.DefaultOffsetUpdateIntervalMs,
			RetryDetails: RetryConfigurations{
				TimeoutMs: base.DefaultOffsetRetryTimeOutMs,
				Count:     base.DefaultOffsetRetryCount,
			},
		},

		LogCollectionDetails: LogCollectionConfigurations{
			CollectionBatchSize:  base.DefaultLogCollectionBatchSize,
			CollectionIntervalMs: base.DefaultLogCollectionIntervalMs,
		},
	}
}

func readConfigYaml() *Config {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	cfg := newConfig()

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}

	checkValuesExistInConfig(cfg)

	reupdateOptionalValuesToDefault(cfg)

	return cfg
}

// checks if all the necessary values that should be configured exist or not
func checkValuesExistInConfig(cfg *Config) {
	switch {
	case len(cfg.LogFiles) == 0:
		panic("Error: log files should be configured for the agent to track and deliver")
	case cfg.AgentLogFile == "":
		panic("Error: agent log file should be configured for log collection")
	case cfg.LogDeliveryDetails.Method == "":
		panic("Error: http method used for log delivery should be configured")
	case cfg.LogDeliveryDetails.Endpoint == "":
		panic("Error: log delivery endpoint should be configured")
	case cfg.LogDeliveryDetails.ExpectedStatusCode == 0:
		panic("Error: expected status code after successful delivery should be given")
	}
}

// reupdates those configurations to default values if any unwelcome value is passed
func reupdateOptionalValuesToDefault(cfg *Config) {
	// handle negative input for log delivery retry configurations
	if cfg.LogDeliveryDetails.RetryDetails.Count < 0 {
		cfg.LogDeliveryDetails.RetryDetails.Count = base.DefaultDeliveryRetryCount
	}
	if cfg.LogDeliveryDetails.RetryDetails.TimeoutMs < 0 {
		cfg.LogDeliveryDetails.RetryDetails.TimeoutMs = base.DefaultDeliveryRetryTimeoutMs
	}

	// handle negative input for log collection configurations
	if cfg.LogCollectionDetails.CollectionBatchSize < 0 {
		cfg.LogCollectionDetails.CollectionBatchSize = base.DefaultLogCollectionBatchSize
	}
	if cfg.LogCollectionDetails.CollectionIntervalMs < 0 {
		cfg.LogCollectionDetails.CollectionIntervalMs = base.DefaultLogCollectionIntervalMs
	}

	// handle negative input for offset persistence configurations
	if cfg.OffsetPersistenceDetails.UpdateIntervalMs < 0 {
		cfg.OffsetPersistenceDetails.UpdateIntervalMs = base.DefaultOffsetUpdateIntervalMs
	}
	if cfg.OffsetPersistenceDetails.RetryDetails.Count < 0 {
		cfg.OffsetPersistenceDetails.RetryDetails.Count = base.DefaultOffsetRetryCount
	}
	if cfg.OffsetPersistenceDetails.RetryDetails.TimeoutMs < 0 {
		cfg.OffsetPersistenceDetails.RetryDetails.TimeoutMs = base.DefaultOffsetRetryTimeOutMs
	}
}

var cfg *Config

func init() {
	cfg = readConfigYaml()
}

func Get() *Config {
	return cfg
}
