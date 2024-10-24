package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/handlers/loans"
	"github.com/wsugiri/loansystem/handlers/master"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/api/loans", loans.ProposeLoan)
	app.Put("/api/loans/:id/approve", loans.ApproveLoan)
	app.Put("/api/loans/:id/reject", loans.RejectLoan)
	app.Put("/api/loans/:id/invest", loans.InvestLoan)
	app.Put("/api/loans/:id/disburse", loans.DisburseLoan)
	app.Get("/api/loans", loans.ListLoan)
	app.Get("/api/loans/:id/schedule", loans.GetScheduleLoan)

	app.Get("/api/master/users", master.ListAllUser)
	app.Get("/api/master/users/:role", master.ListAllUser)
}
