package entity

import "time"

// The log format to be emitted by the services
type Log struct {
	Level       string                 `json:"level"`              // represents the log level
	Timestamp   time.Time              `json:"timestamp"`          // represents the log creation time
	ServiceName string                 `json:"service"`            // represents the service that emitted the log
	Metadata    map[string]interface{} `json:"metadata,omitempty"` // holds extra log context
	Message     string                 `json:"message"`            // the log message
}
