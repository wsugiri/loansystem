package loans

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func DisburseLoan(c *fiber.Ctx) error {
	var payload struct {
		InvestorID int     `json:"investor_id"`
		Amount     float32 `json:"amount"`
	}
	loanId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Check Investor
	var query string
	var investor models.User

	query = `select id, name, role from users where id = ?`
	if err := utils.DB.QueryRow(query, payload.InvestorID).Scan(&investor.ID, &investor.Name, &investor.Role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if investor.Role != "investor" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid_investor",
		})
	}

	fmt.Println(loanId)

	queries := c.Queries()

	fmt.Println(queries)
	fmt.Println(queries["brand"])

	return c.JSON(queries)
}
