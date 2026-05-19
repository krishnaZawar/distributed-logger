package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krishnaZawar/distributed-logger/ingestion-node/internal/router"
)

func main() {
	app := fiber.New()

	router.New(app)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
