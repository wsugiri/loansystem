package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/wsugiri/loansystem/routers"
	"github.com/wsugiri/loansystem/utils"
)

func SetupAndRunApp() error {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	utils.InitializeDatabase()
	defer utils.DB.Close() // Close the DB connection on exit

	// Initialize Fiber app
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Sample 001")
	})

	// Initialize Routings
	routers.SetupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "9001"
	}
	app.Listen(":" + port)

	return nil
}
