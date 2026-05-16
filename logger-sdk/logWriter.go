package logger

import (
	"encoding/json"
	"io"
	"sync"
)

// Will be used internally by the Log to append the log to the given file
type logWriter struct {
	// mu protects the writer from concurrent access during writes
	mu     sync.RWMutex
	writer io.Writer
}

func newLogWriter(writer io.Writer) *logWriter {
	return &logWriter{writer: writer}
}

// writes the log to the writer pointed by the logWriter
func (lw *logWriter) WriteLog(log *log) {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	data, _ := json.Marshal(log)
	data = append(data, '\n')
	_, _ = lw.writer.Write(data)
}
