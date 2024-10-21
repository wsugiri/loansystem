package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/handlers/loans"
	"github.com/wsugiri/loansystem/handlers/master"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/api/loans", loans.ListLoan)
	app.Post("/api/loans", loans.CreateLoan)
	app.Put("/api/loans/:id/approve", loans.ApproveLoan)
	app.Put("/api/loans/:id/reject", loans.RejectLoan)

	app.Get("/api/master/db", master.GetDBConnection)
	app.Get("/api/master/users", master.ListAllUser)
	app.Get("/api/master/users/:role", master.ListAllUser)
}
