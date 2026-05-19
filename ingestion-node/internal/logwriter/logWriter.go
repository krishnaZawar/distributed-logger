package logwriter

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/krishnaZawar/distributed-logger/ingestion-node/internal/config"
	"github.com/krishnaZawar/distributed-logger/ingestion-node/internal/entity"
)

// Will be used internally by the node to persist the logs to the given file
type logWriter struct {
	// mu protects the writer from concurrent access during writes
	mu     sync.RWMutex
	writer io.Writer
}

func newLogWriter(writer io.Writer) *logWriter {
	return &logWriter{writer: writer}
}

// writes the log to the writer pointed by the logWriter
//
// returns the logs written and error if any
func (lw *logWriter) WriteLogs(logs []entity.Log) (int, error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	for idx := range logs {
		data, _ := json.Marshal(logs[idx])
		data = append(data, '\n')
		_, err := lw.writer.Write(data)
		if err != nil {
			return idx, err
		}
	}
	return len(logs), nil
}

var writer *logWriter

func init() {
	cfg := config.Get()
	file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	writer = newLogWriter(file)
}

// This is function writes logs to the configured file for persistence
//
// On failure or incomplete writes, it retries once for log persistence
func WriteLogs(logs []entity.Log) {
	written, err := writer.WriteLogs(logs)
	if err != nil || written < len(logs) {
		writer.WriteLogs(logs[written:])
	}
}
