package logger

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer

	service := "test-srv"

	logger := New(service, &buf)

	if logger == nil {
		t.Fatal("expected logger to be initialized")
	}

	if logger.service != service {
		t.Fatalf("expected service to be %s, got %s", service, logger.service)
	}

	if logger.writer == nil {
		t.Fatal("expected writer to be initialized")
	}
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer

	service := "test-srv"

	logger := New(service, &buf)

	tests := []struct {
		name  string
		event *LogEvent
		level logLevel
	}{
		{
			name:  "debug logger",
			event: logger.Debug(),
			level: LevelDebug,
		},
		{
			name:  "info logger",
			event: logger.Info(),
			level: LevelInfo,
		},
		{
			name:  "warn logger",
			event: logger.Warn(),
			level: LevelWarn,
		},
		{
			name:  "error logger",
			event: logger.Error(),
			level: LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.event == nil {
				t.Fatal("expected log event to be initialized")
			}

			if tt.event.log.Level != tt.level {
				t.Fatalf("expected level %v, got %v", tt.level, tt.event.log.Level)
			}

			if tt.event.log.ServiceName != service {
				t.Fatalf("expected service %s, got %s", service, tt.event.log.ServiceName)
			}

			if tt.event.writer == nil {
				t.Fatal("expected writer to be initialized")
			}
		})
	}
}

func TestLoggerThreadSafety(t *testing.T) {
	var buf bytes.Buffer

	service := "concurreny-test-srv"

	logger := New(service, &buf)

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			_ = logger.Debug()
			_ = logger.Info()
			_ = logger.Warn()
			_ = logger.Error()
		}()
	}

	wg.Wait()
}

func TestLoggerUsesProvidedWriter(t *testing.T) {
	var buf bytes.Buffer

	service := "test-writer-srv"

	logger := New(service, &buf)

	event := logger.Info()

	if event.writer == nil {
		t.Fatal("expected writer to exist")
	}

	// sanity check that underlying writer is functional
	_, err := event.writer.writer.Write([]byte("test log"))
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	if !strings.Contains(buf.String(), "test log") {
		t.Fatal("expected data to be written to buffer")
	}
}
