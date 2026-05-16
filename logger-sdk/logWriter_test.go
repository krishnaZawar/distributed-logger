package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestWriteLog_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	w := newLogWriter(&buf)

	service := "test-srv"
	message := "test message"

	l := &log{
		Level:       LevelInfo,
		Timestamp:   time.Now(),
		ServiceName: service,
		Message:     message,
	}

	w.WriteLog(l)

	out := strings.TrimSpace(buf.String())

	var decoded log
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}

	if decoded.Level != LevelInfo {
		t.Fatalf("expected level %s got %s", LevelInfo, decoded.Level)
	}

	if decoded.ServiceName != service {
		t.Fatalf("expected service %s got %s", service, decoded.ServiceName)
	}

	if decoded.Message != message {
		t.Fatalf("expected message %s got %s", message, decoded.Message)
	}
}

func TestWriteLog_IsThreadSafe(t *testing.T) {
	var buf bytes.Buffer
	w := newLogWriter(&buf)

	service := "test-srv"
	message := "test message"

	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func(i int) {
			w.WriteLog(&log{
				Level:       LevelInfo,
				Timestamp:   time.Now(),
				ServiceName: service,
				Message:     message,
			})
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < 50; i++ {
		<-done
	}
}
