package logger

import (
	"io"
)

// will be used across services for standardized logging
//
// Used to enforce strict function use for logging. This is thread safe
type Logger struct {
	// The logger is thread safe because the fields are immutable after initialization making it read safe by default
	// hence no usage of mutex here

	service string     // service instantiating the logger
	writer  *logWriter // writes the logEvent to the file
}

// creates a new thread safe logger instance for the service to use
func New(service string, writer io.Writer) *Logger {
	return &Logger{
		service: service,
		writer:  newLogWriter(writer),
	}
}

// initiates a debug level log
func (logger *Logger) Debug() *LogEvent {
	return newLogEvent(LevelDebug, logger.service, logger.writer)
}

// initiates a info level log
func (logger *Logger) Info() *LogEvent {
	return newLogEvent(LevelInfo, logger.service, logger.writer)
}

// initiates a warn level log
func (logger *Logger) Warn() *LogEvent {
	return newLogEvent(LevelWarn, logger.service, logger.writer)
}

// initiates a error level log
func (logger *Logger) Error() *LogEvent {
	return newLogEvent(LevelError, logger.service, logger.writer)
}
