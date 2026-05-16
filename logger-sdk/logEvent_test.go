package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestNewLogEvent(t *testing.T) {
	var buf bytes.Buffer
	writer := newLogWriter(&buf)

	service := "test-srv"

	event := newLogEvent(LevelInfo, service, writer)

	if event == nil {
		t.Fatal("expected log event to be initialized")
	}

	if event.log == nil {
		t.Fatal("expected log to be initialized")
	}

	if event.log.Level != LevelInfo {
		t.Fatalf("expected level %s, got %s", LevelInfo, event.log.Level)
	}

	if event.log.ServiceName != service {
		t.Fatalf("expected service %s, got %s", service, event.log.ServiceName)
	}

	if event.writer == nil {
		t.Fatal("expected writer to be initialized")
	}

	if event.log.Timestamp.Sub(time.Now()) > 0 {
		t.Fatal("expected timestamp to be recent")
	}
}

func TestLogEventMsg(t *testing.T) {
	var buf bytes.Buffer
	writer := newLogWriter(&buf)

	service := "test-srv"

	event := newLogEvent(LevelDebug, service, writer)

	message := "user logged in"
	event.Msg(message)

	if event.log.Message != message {
		t.Fatalf("expected message '%s', got %s", message, event.log.Message)
	}

	output := buf.String()

	if !strings.Contains(output, message) {
		t.Fatal("expected log output to contain message")
	}

	if !strings.Contains(output, service) {
		t.Fatal("expected log output to contain service name")
	}

	if !strings.Contains(output, string(LevelDebug)) {
		t.Fatal("expected log output to contain log level")
	}
}

func TestLogEventMsgf(t *testing.T) {
	var buf bytes.Buffer
	writer := newLogWriter(&buf)

	service := "test-srv"

	event := newLogEvent(LevelError, service, writer)

	event.Msgf("payment failed for user %d", 42)

	expected := "payment failed for user 42"

	if event.log.Message != expected {
		t.Fatalf("expected message %q, got %q", expected, event.log.Message)
	}

	output := buf.String()

	if !strings.Contains(output, expected) {
		t.Fatal("expected formatted message in log output")
	}

	if !strings.Contains(output, service) {
		t.Fatal("expected service name in log output")
	}

	if !strings.Contains(output, string(LevelError)) {
		t.Fatal("expected log level in log output")
	}
}

func TestLogSerialization(t *testing.T) {
	logEntry := &log{
		Level:       LevelWarn,
		Timestamp:   time.Now(),
		ServiceName: "test-srv",
		Message:     "test message",
	}

	data, err := json.Marshal(logEntry)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	output := string(data)

	if !strings.Contains(output, `"level":"warn"`) {
		t.Fatal("expected serialized level field")
	}

	if !strings.Contains(output, `"service":"test-srv"`) {
		t.Fatal("expected serialized service field")
	}

	if !strings.Contains(output, `"message":"test message"`) {
		t.Fatal("expected serialized message field")
	}
}

func TestLogEvent_WithMetadata(t *testing.T) {
	event := &LogEvent{
		log: &log{
			Metadata: make(map[string]interface{}),
		},
	}

	key := "user_id"
	val := 123

	result := event.WithMetadata(key, val)

	// assert return value (method chaining)
	if result != event {
		t.Errorf("expected same LogEvent pointer to be returned")
	}

	// assert metadata stored correctly
	if event.log.Metadata[key] != val {
		t.Errorf("expected metadata %s=%v, got %v", key, val, event.log.Metadata["user_id"])
	}
}
