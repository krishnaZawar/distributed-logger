package logger

import (
	"encoding/json"
	"io"
)

// Will be used internally by the Log to append the log to the given file
type logWriter struct {
	writer io.Writer
}

func newLogWriter(writer io.Writer) *logWriter {
	return &logWriter{writer: writer}
}

// writes the log to the writer pointed by the logWriter
func (lw *logWriter) WriteLog(log *log) {
	data, _ := json.Marshal(log)
	data = append(data, '\n')
	_, _ = lw.writer.Write(data)
}
