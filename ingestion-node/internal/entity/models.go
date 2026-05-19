package entity

import "time"

type Log struct {
	// The data sent by the logging agents
	Level       string                 `json:"level"`              // represents the log level
	Timestamp   time.Time              `json:"timestamp"`          // represents the log creation time
	ServiceName string                 `json:"service"`            // represents the service that emitted the log
	Metadata    map[string]interface{} `json:"metadata,omitempty"` // holds extra log context
	Message     string                 `json:"message"`            // the log message

	// The data that should be enriched by the ingestion node
	IP string `json:"host"` // the host, here the logging agent, sending the logs
}
