package loans

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func ListLoan(c *fiber.Ctx) error {
	var (
		query string
		rows  *sql.Rows
		err   error
	)

	query = `select id,borrower_id, principal_amount, rate, duration_weeks, status from loans`

	if c.Query("status") != "" {
		rows, err = utils.DB.Query(query+" where status = ?", c.Query("status"))
	} else {
		rows, err = utils.DB.Query(query)
	}

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	defer rows.Close()

	var loans []models.Loan
	for rows.Next() {
		var loan models.Loan
		if err := rows.Scan(&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.Rate, &loan.DurationWeek, &loan.Status); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}
		loans = append(loans, loan)
	}

	response := fiber.Map{
		"status": "success",
		"total":  len(loans),
		"items":  loans,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
