package main

import (
	"github.com/krishnaZawar/distributed-logger/logging-agent/internal/entity"
)

func main() {
	agent := entity.NewLoggingAgent()
	agent.Read()
}
