package loans

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func ProposeLoan(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	query = `select id, name, role from users where id = ? and role = 'borrower'`
	if err := utils.DB.QueryRow(query, payload.BorrowerID).Scan(&user.ID, &user.Name, &user.Role); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid borrower_id",
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Calculate total loan and instalment
	totalLoan := payload.PrincipalAmount + (payload.PrincipalAmount * payload.InterestRate * 0.01)
	instalment := totalLoan / payload.LoanDurationWeeks
	status := "proposed"

	query = `insert into loans (borrower_id, principal_amount, rate, total_loan, instalment, duration_weeks, status, agreement_url) values (?, ?, ?, ?, ?, ?, ?, ?)`

	insertResult, err := utils.DB.Exec(query,
		payload.BorrowerID, payload.PrincipalAmount, payload.InterestRate,
		totalLoan, instalment, payload.LoanDurationWeeks, status, payload.AgreementUrl,
	)

	if err != nil {
		panic(err)
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		panic(err)
	}

	response = fiber.Map{
		"status":  "success",
		"message": "Loan successfully created",
		"data": fiber.Map{
			"loan_id":             id,
			"borrower_id":         payload.BorrowerID,
			"principal_amount":    payload.PrincipalAmount,
			"interest_rate":       payload.InterestRate,
			"loan_duration_weeks": payload.LoanDurationWeeks,
			"total_loan":          totalLoan,
			"agreement_url":       payload.AgreementUrl,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
