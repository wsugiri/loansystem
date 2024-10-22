package loans

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func InvestLoan(c *fiber.Ctx) error {
	var payload struct {
		InvestorID int     `json:"investor_id"`
		Amount     float32 `json:"amount"`
	}
	loanId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Check Investor
	var query string
	var investor models.User

	query = `select id, name, email, role from users where id = ?`
	if err := utils.DB.QueryRow(query, payload.InvestorID).Scan(&investor.ID, &investor.Name, &investor.Email, &investor.Role); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invalid investor id"})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if investor.Role != "investor" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid investor id",
		})
	}

	// Check Loan
	var loan models.Loan
	query = `
	select a.id, a.borrower_id, a.principal_amount, a.status
	     , sum(ifnull(b.amount, 0)) as invested_amount
	  from loans a
	  left join investments b on b.loan_id = a.id 
	 where a.id = ?
	   and a.status in ('approved', 'invested')
	 group by a.id`

	if err := utils.DB.QueryRow(query, loanId).Scan(&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.Status, &loan.InvestedAmount); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invalid loan id"})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	availableAmount := loan.PrincipalAmount - loan.InvestedAmount

	if availableAmount < payload.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("cannot invest more than %d", int(availableAmount))})
	}

	// Insert Investment
	query = "insert into investments (loan_id, investor_id, amount) values (?, ?, ?)"
	_, err = utils.DB.Exec(query, loanId, payload.InvestorID, payload.Amount)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Update Status
	status := "invested"
	query = "update loans set status = ? where id = ?"
	_, err = utils.DB.Exec(query, status, loanId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"prncipal_amount":  loan.PrincipalAmount,
			"invested_amount":  loan.InvestedAmount + payload.Amount,
			"available_amount": loan.PrincipalAmount - (loan.InvestedAmount + payload.Amount),
		},
		"investor": fiber.Map{
			"name":  investor.Name,
			"email": investor.Email,
		},
	})
}
