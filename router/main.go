package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/handlers/loans"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/api/loans", loans.ListLoan)
	app.Post("/api/loans", loans.CreateLoan)
	app.Put("/api/loans/:id/approve", loans.ApproveLoan)
	app.Put("/api/loans/:id/reject", loans.RejectLoan)
}
