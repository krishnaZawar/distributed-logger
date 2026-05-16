# Logger SDK

A lightweight, thread-safe structured logging package for Go services.

It provides leveled logging (`debug`, `info`, `warn`, `error`) with JSON output and a simple fluent API designed for standardized logging across multiple services.

---

## Features

- Structured JSON logs
- Log levels
    - debug
    - info
    - warn
    - error
- Thread-safe concurrent writes
- Simple fluent API
- Service-based log tagging

---

## Installation

```bash
go get github.com/krishnaZawar/distributed-logger/logger
```

## Quick Start

```go
package main

import (
	"os"

	"github.com/krishnaZawar/distributed-logger/logger"
)

func main() {
	log := logger.New("test-srv", os.Stdout)

	log.Info().Msg("service started")
}
```

## Usage

### Logger creation

The logger is bound to a service name which is included in every log entry.

```go
logger := logger.New("test-service", os.Stdout)
```

### Log Levels

```go
log.Debug().Msg("debug message")                    // debug level log
log.Info().Msg("user logged in")                    // info level log
log.Warn().Msg("slow response detected")            // warn level log
log.Error().Msg("database connection failed")       // error level log
```

### Formatted logging

You can format log messages using `Msgf`

```go
log.Info().Msgf("user %d logged in", userID)
```

### Log Structure

All logs are emitted in JSON format

```json
{
  "level": "info",
  "timestamp": "2026-05-16T10:00:00Z",
  "service": "auth-service",
  "message": "user logged in"
}
```

## Future Use

1. Allow logging of extra metadata with the log