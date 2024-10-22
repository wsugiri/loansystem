package loans

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func CreateLoan(c *fiber.Ctx) error {
	// Define the variables
	var query string
	var user models.User
	var payload struct {
		BorrowerID        int     `json:"borrower_id"`
		PrincipalAmount   float32 `json:"principal_amount"`
		InterestRate      float32 `json:"interest_rate"`
		LoanDurationWeeks float32 `json:"loan_duration_weeks"`
		AgreementUrl      string  `json:"agreement_url"`
	}
	var response fiber.Map

	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query = `select id, name, role from users where id = ?`
	err := utils.DB.QueryRow(query, payload.BorrowerID).Scan(&user.ID, &user.Name, &user.Role)

	if err != nil {
		panic(err)
	}

	if user.Role != "borrower" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "unregistered borrower",
		})
	}

	// Calculate total loan and instalment
	totalLoan := payload.PrincipalAmount + (payload.PrincipalAmount * payload.InterestRate * 0.01)
	instalment := totalLoan / payload.LoanDurationWeeks
	status := "proposed"

	query = `insert into loans 
	(borrower_id, principal_amount, rate, total_loan, instalment, status, agreement_url) values (?, ?, ?, ?, ?, ?, ?)`

	insertResult, err := utils.DB.Exec(query,
		payload.BorrowerID, payload.PrincipalAmount, payload.InterestRate,
		totalLoan, instalment, status, payload.AgreementUrl,
	)

	if err != nil {
		panic(err)
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		panic(err)
	}

	response = fiber.Map{
		"id":     id,
		"status": status,
		"data": fiber.Map{
			"total_loan":     totalLoan,
			"instalment":     instalment,
			"duration_weeks": payload.LoanDurationWeeks,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
