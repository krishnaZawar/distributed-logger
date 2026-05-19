package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/distributed-logger/ingestion-node/internal/entity"
	"github.com/krishnaZawar/distributed-logger/ingestion-node/internal/logwriter"
)

// Ingests logs from various logging agents
//
// Also processes and enriches the logs and sends them for downstream processing
// Here just writing to a file
func IngestLogs(ctx *fiber.Ctx) error {
	var logs []entity.Log
	if err := ctx.BodyParser(&logs); err != nil {
		return ctx.Status(400).SendString(fmt.Sprintf("Bad Request: %s", err.Error()))
	}

	// enrich the logs with IP
	for idx := range logs {
		logs[idx].IP = ctx.IP()
	}

	// write to log file
	logwriter.WriteLogs(logs)

	return ctx.Status(200).SendString("logs delivered")
}

func Ping(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(map[string]string{
		"data": "pong.",
	})
}
