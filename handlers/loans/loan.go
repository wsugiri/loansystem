package loans

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func CreateLoan(c *fiber.Ctx) error {
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
	if err := utils.DB.QueryRow(query, payload.BorrowerID).Scan(&user.ID, &user.Name, &user.Role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if user.Role != "borrower" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "unregistered_borrower",
		})
	}

	// Calculate total loan and instalment
	totalLoan := payload.PrincipalAmount + (payload.PrincipalAmount * payload.InterestRate * 0.01)
	instalment := totalLoan / payload.LoanDurationWeeks
	status := "proposed"

	query = `insert into loans 
	(borrower_id, principal_amount, rate, total_loan, instalment, duration_weeks, status, agreement_url) values (?, ?, ?, ?, ?, ?, ?, ?)`

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

func ApproveLoan(c *fiber.Ctx) error {
	var payload struct {
		EmployeeID     int    `json:"employee_id"`
		ApprovalDate   string `json:"approval_date"`
		ValidatorPhoto string `json:"validator_photo"`
	}
	loanId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		panic(err)
	}

	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	status := "approved"
	var query string
	var user models.User

	query = `select id, name, role from users where id = ?`
	if err := utils.DB.QueryRow(query, payload.EmployeeID).Scan(&user.ID, &user.Name, &user.Role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if user.Role != "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "unregistered_approver",
		})
	}

	var loan models.Loan

	query = `select id, total_loan, instalment, duration_weeks, status from loans where id = ?`
	if err := utils.DB.QueryRow(query, loanId).Scan(&loan.ID, &loan.TotalLoan, &loan.Instalment, &loan.DurationWeek, &loan.Status); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invalid loan id"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if loan.Status == "approved" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "status_already_approved",
		})
	}

	if loan.Status != "proposed" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_status",
		})
	}

	// Parse the approve date.
	approveDate, err := time.Parse("2006-01-02", payload.ApprovalDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid approval date"})
	}

	// Generate the instalments
	instalments := make([]models.Instalment, 0, loan.DurationWeek)
	for idx := 0; idx < loan.DurationWeek; idx++ {
		// Calculate the due date by adding (idx * 7) days.
		dueDate := approveDate.AddDate(0, 0, (idx+1)*7).Format("2006-01-02")

		instalments = append(instalments, models.Instalment{
			Week:    idx + 1,
			Amount:  loan.Instalment,
			DueDate: dueDate,
		})
	}

	// Inser installments
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

	query = "update loans set status = ? where id = ?"
	_, err = utils.DB.Exec(query, status, loanId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp := fiber.Map{
		"id":     loanId,
		"status": status,
		"datas": fiber.Map{
			"duration_weeks": loan.DurationWeek,
			"total_loan":     loan.TotalLoan,
			"instalments":    instalments,
		},
	}

	return c.JSON(resp)
}

func RejectLoan(c *fiber.Ctx) error {
	var payload struct {
		EmployeeID       int    `json:"employee_id"`
		RejectionDate    string `json:"rejection_date"`
		RejectionMessage string `json:"rejection_message"`
	}
	loanId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		panic(err)
	}

	// Parse the request body
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	status := "rejected"
	var query string
	var user models.User

	query = `select id, name, role from users where id = ?`
	if err := utils.DB.QueryRow(query, payload.EmployeeID).Scan(&user.ID, &user.Name, &user.Role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if user.Role != "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "unregistered_rejector",
		})
	}

	var loan models.Loan

	query = `select id, total_loan, instalment, duration_weeks, status from loans where id = ?`
	if err := utils.DB.QueryRow(query, loanId).Scan(&loan.ID, &loan.TotalLoan, &loan.Instalment, &loan.DurationWeek, &loan.Status); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invalid loan id"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if loan.Status != "proposed" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_status",
		})
	}

	// Parse the approve date.
	rejectionDate, err := time.Parse("2006-01-02", payload.RejectionDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid approval date"})
	}

	fmt.Println(rejectionDate)

	// // Generate the instalments
	// instalments := make([]models.Instalment, 0, loan.DurationWeek)
	// for idx := 0; idx < loan.DurationWeek; idx++ {
	// 	// Calculate the due date by adding (idx * 7) days.
	// 	dueDate := approveDate.AddDate(0, 0, (idx+1)*7).Format("2006-01-02")

	// 	instalments = append(instalments, models.Instalment{
	// 		Week:    idx + 1,
	// 		Amount:  loan.Instalment,
	// 		DueDate: dueDate,
	// 	})
	// }

	// // Inser installments
	// query = "insert into payments (loan_id, week, amount, due_date) values "

	// // Create placeholders for each row.
	// values := []string{}
	// args := []interface{}{}

	// for _, instalment := range instalments {
	// 	values = append(values, "(?, ?, ?, ?)")
	// 	args = append(args, loanId, instalment.Week, instalment.Amount, instalment.DueDate)
	// }

	// // Join all value placeholders into the query.
	// query += strings.Join(values, ", ")

	// // Execute the query.
	// _, err = utils.DB.Exec(query, args...)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }

	// query = "update loans set status = ? where id = ?"
	// _, err = utils.DB.Exec(query, status, loanId)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }

	resp := fiber.Map{
		"id":             loanId,
		"status":         status,
		"rejection_date": payload.RejectionDate,
	}

	return c.JSON(resp)
}
