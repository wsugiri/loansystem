package loans

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/handlers/loans/constants"
	"github.com/wsugiri/loansystem/models"
)

func GetOutstanding(c *fiber.Ctx) error {
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

	var transDate = c.Query("trans_date")
	if transDate == "" {
		transDate = time.Now().Format("2006-01-02")
	}

	println(transDate)

	loan, _ := CheckLoanOutstanding(loanId, transDate)

	return c.JSON(models.Response{
		Status:  "success",
		Message: "Outstanding balance retrieved successfully",
		Data: fiber.Map{
			"loan_id":            loanId,
			"outstanding_amount": loan.OutstandingLoan,
		},
	})
}
