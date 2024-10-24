package loans

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/handlers/loans/constants"
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func DisburseLoan(c *fiber.Ctx) error {
	var payload struct {
		EmployeeID       int    `json:"employee_id"`
		DisbursementDate string `json:"disbursement_date"`
		AgreementLetter  string `json:"agreement_letter"`
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

	// Check Employee
	var query string
	var employee models.User

	fmt.Println(payload)

	query = `select id, name, email, role from users where id = ? and role = 'staff'`
	if err := utils.DB.QueryRow(query, payload.EmployeeID).Scan(&employee.ID, &employee.Name, &employee.Email, &employee.Role); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": constants.ErrEmployeeInvalid,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	// Check Loan
	var loan models.Loan
	query = `
	select a.id, a.borrower_id, a.principal_amount, a.status, a.rate, a.duration_weeks, instalment
	     , sum(ifnull(b.amount, 0)) as invested_amount
	  from loans a
	  left join investments b on b.loan_id = a.id 
	 where a.id = ?
	   and a.status in ('approved', 'invested','disbursed')
	 group by a.id`

	if err := utils.DB.QueryRow(query, loanId).Scan(&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.Status, &loan.Rate, &loan.DurationWeek, &loan.Instalment, &loan.InvestedAmount); err != nil {
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

	if loan.Status == "disbursed" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": constants.ErrLoanDisbursed,
		})
	}

	// Insert Disbursement
	query = "insert into disbursements (loan_id, field_officer_id, disbursement_date, agreement_signed_url) values (?, ?, ?, ?)"
	_, err = utils.DB.Exec(query, loanId, payload.EmployeeID, payload.DisbursementDate, payload.AgreementLetter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Update Status
	status := "disbursed"
	query = "update loans set status = ? where id = ?"
	_, err = utils.DB.Exec(query, status, loanId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate the instalments
	disbursedDate, err := time.Parse("2006-01-02", payload.DisbursementDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid approval date"})
	}

	instalments := make([]models.Instalment, 0, loan.DurationWeek)
	for idx := 0; idx < loan.DurationWeek; idx++ {
		// Calculate the due date by adding (idx * 7) days.
		dueDate := disbursedDate.AddDate(0, 0, (idx+1)*7).Format("2006-01-02")

		instalments = append(instalments, models.Instalment{
			Week:    idx + 1,
			Amount:  loan.Instalment,
			DueDate: dueDate,
		})
	}

	// Insert installments
	query = "insert into payments (loan_id, week, amount, due_date) values "

	// Create placeholders for each row.
	values := []string{}
	args := []interface{}{}

	for _, instalment := range instalments {
		values = append(values, "(?, ?, ?, ?)")
		args = append(args, loanId, instalment.Week, instalment.Amount, instalment.DueDate)
	}

	// Join all value placeholders into the query.
	query += strings.Join(values, ", ")

	// Execute the query.
	_, err = utils.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Loan successfully disbursed",
		"data": fiber.Map{
			"loan_id":           loanId,
			"employee_id":       employee.ID,
			"disbursement_date": payload.DisbursementDate,
			"agreement_letter":  payload.AgreementLetter,
			"disbursed_amount":  loan.PrincipalAmount,
			"borrower_id":       loan.BorrowerID,
		},
	})
}
