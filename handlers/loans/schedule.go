package loans

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/handlers/loans/constants"
	"github.com/wsugiri/loansystem/models"
)

func GetScheduleLoan(c *fiber.Ctx) error {
	loanId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	_, err = CheckLoan(loanId)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": constants.ErrLoanInvalid,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	schedules, _ := GetPayments(loanId)

	return c.JSON(models.Response{
		Status:  "success",
		Message: "Loan schedule retrieved successfully.",
		Data: fiber.Map{
			"loan_id":  loanId,
			"schedule": schedules,
		},
	})
}
