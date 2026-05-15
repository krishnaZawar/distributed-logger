package logger

import "io"

// will be used across services for standardized logging
//
// Used to enforce strict function use for logging
type Logger struct {
	Service string // service instantiating the logger
	writer  *logWriter
}

// creates a new logger instance for the service to use
func New(service string, writer io.Writer) *Logger {
	return &Logger{
		Service: service,
		writer:  newLogWriter(writer),
	}
}

// initiates a debug level log
func (logger *Logger) Debug() *LogEvent {
	return newLogEvent(LevelDebug, logger.Service, logger.writer)
}

// initiates a info level log
func (logger *Logger) Info() *LogEvent {
	return newLogEvent(LevelInfo, logger.Service, logger.writer)
}

// initiates a warn level log
func (logger *Logger) Warn() *LogEvent {
	return newLogEvent(LevelWarn, logger.Service, logger.writer)
}

// initiates a error level log
func (logger *Logger) Error() *LogEvent {
	return newLogEvent(LevelError, logger.Service, logger.writer)
}
