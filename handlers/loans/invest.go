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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	// Check Investor
	var query string
	var investor models.User

	query = `select id, name, email, role from users where id = ? and role = 'investor'`
	if err := utils.DB.QueryRow(query, payload.InvestorID).Scan(&investor.ID, &investor.Name, &investor.Email, &investor.Role); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "The provided investor_id is invalid or does not exist. Please check and try again",
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Check Loan
	var loan models.Loan
	query = `
	select a.id, a.borrower_id, a.principal_amount, a.status, a.rate
	     , sum(ifnull(b.amount, 0)) as invested_amount
	  from loans a
	  left join investments b on b.loan_id = a.id 
	 where a.id = ?
	   and a.status in ('approved', 'invested')
	 group by a.id`

	if err := utils.DB.QueryRow(query, loanId).Scan(&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.Status, &loan.Rate, &loan.InvestedAmount); err != nil {
		if err.Error() == "sql: no rows in result set" {
			if err.Error() == "sql: no rows in result set" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status":  "error",
					"message": "The provided loan_id is invalid or does not exist. Please check and try again",
				})
			}
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	availableAmount := loan.PrincipalAmount - loan.InvestedAmount

	if availableAmount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "The loan has already been fully funded and cannot accept further investments.",
		})
	}

	if availableAmount < payload.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("Investment exceeds the maximum allowable amount of " + utils.InsertCommas(int(availableAmount)) + ". Please adjust the amount and try again"),
		})
	}

	// Insert Investment
	query = "insert into investments (loan_id, investor_id, amount) values (?, ?, ?)"
	_, err = utils.DB.Exec(query, loanId, payload.InvestorID, payload.Amount)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Update Status
	status := "invested"
	query = "update loans set status = ? where id = ?"
	_, err = utils.DB.Exec(query, status, loanId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Investment successfully made",
		"data": fiber.Map{
			"loan_id":               loanId,
			"borrower_id":           loan.BorrowerID,
			"investor_id":           payload.InvestorID,
			"investment_amount":     payload.Amount,
			"total_invested":        loan.InvestedAmount + payload.Amount, // Total amount invested after this investment
			"loan_principal_amount": loan.PrincipalAmount,
			"remaining_amount":      loan.PrincipalAmount - (loan.InvestedAmount + payload.Amount), // Amount still needed to fully fund the loan
		},
	})
}
