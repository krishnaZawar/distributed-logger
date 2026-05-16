package logger

import (
	"fmt"
	"time"
)

// indicates the level of the log
//
// 4 levels defined: debug, info, warn, error
type logLevel string

const (
	LevelDebug logLevel = "debug"
	LevelInfo  logLevel = "info"
	LevelWarn  logLevel = "warn"
	LevelError logLevel = "error"
)

// The log format to be emitted by the services
type log struct {
	Level       logLevel  `json:"level"`     // represents the log level
	Timestamp   time.Time `json:"timestamp"` // represents the log creation time
	ServiceName string    `json:"service"`   // represents the service that emitted the log
	Message     string    `json:"message"`   // the log message
}

// The event emitted by the logger functions
//
// holds context related to the current log emitted
//
// can be extended accordingly for specific use
type LogEvent struct {
	log    *log
	writer *logWriter
}

func newLogEvent(level logLevel, service string, writer *logWriter) *LogEvent {
	return &LogEvent{
		log: &log{
			Level:       level,
			Timestamp:   time.Now(),
			ServiceName: service,
		},
		writer: writer,
	}
}

// assigns message to the Log
func (event *LogEvent) Msg(message string) {
	event.log.Message = message

	// push log to file
	event.writer.WriteLog(event.log)
}

// assigns formatted message to the Log
func (event *LogEvent) Msgf(format string, a ...any) {
	event.log.Message = fmt.Sprintf(format, a...)

	// push log to file
	event.writer.WriteLog(event.log)
}
