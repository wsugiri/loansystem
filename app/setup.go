package app

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/router"
)

func SetupAndRunApp() error {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Sample 001")
	})

	router.SetupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9001"
	}
	app.Listen(":" + port)

	return nil
}
