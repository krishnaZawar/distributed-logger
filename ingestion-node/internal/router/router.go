package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/distributed-logger/ingestion-node/internal/handler"
)

func New(app *fiber.App) {
	app.Get("/ping", handler.Ping)

	app.Post("/ingest", handler.IngestLogs)
}
