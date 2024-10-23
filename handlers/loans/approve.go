package loans

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	status := "approved"
	var query string
	var user models.User

	query = `select id, name, role from users where id = ?`
	if err := utils.DB.QueryRow(query, payload.EmployeeID).Scan(&user.ID, &user.Name, &user.Role); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "The provided employee_id is invalid or does not exist. Please check and try again",
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if user.Role != "staff" {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"status":  "error",
			"message": "The provided employee ID is not authorized to approve this loan",
		})
	}

	var loan models.Loan

	query = `select id, total_loan, instalment, duration_weeks, status from loans where id = ?`
	if err := utils.DB.QueryRow(query, loanId).Scan(&loan.ID, &loan.TotalLoan, &loan.Instalment, &loan.DurationWeek, &loan.Status); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid loan id",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if loan.Status == "approved" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Loan status already approved",
		})
	}

	if loan.Status == "rejected" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Loan status already rejected",
		})
	}

	if loan.Status != "proposed" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Loan is not in an proposed state",
		})
	}

	// Insert Approval
	query = "insert into approvals (loan_id, picture_proof_url, approval_date, approval_by) values (?, ?, ?, ?)"
	_, err = utils.DB.Exec(query, loanId, payload.ValidatorPhoto, payload.ApprovalDate, payload.EmployeeID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Update Status
	query = "update loans set status = ? where id = ?"
	_, err = utils.DB.Exec(query, status, loanId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	resp := fiber.Map{
		"status":  "success",
		"message": "Loan successfully approved",
		"datas": fiber.Map{
			"loan_id":         loanId,
			"borrower_id":     loan.BorrowerID,
			"approval_date":   payload.ApprovalDate,
			"employee_id":     payload.EmployeeID,
			"validator_photo": payload.ValidatorPhoto,
			"loan_status":     status,
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

	status := "rejected"
	var query string
	var user models.User

	query = `select id, name, role from users where id = ? and role = 'staff'`
	if err := utils.DB.QueryRow(query, payload.EmployeeID).Scan(&user.ID, &user.Name, &user.Role); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid employee_id",
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	var loan models.Loan

	query = `select id, borrower_id, total_loan, instalment, duration_weeks, status from loans where id = ? and status = 'proposed'`
	if err := utils.DB.QueryRow(query, loanId).Scan(&loan.ID, &loan.BorrowerID, &loan.TotalLoan, &loan.Instalment, &loan.DurationWeek, &loan.Status); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid loan id",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Parse the rejection date.
	rejectionDate, err := time.Parse("2006-01-02", payload.RejectionDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid rejection date",
		})
	}

	// Insert Rejection
	query = "insert into rejections (loan_id, rejection_reason, rejection_date, rejected_by) values (?, ?, ?, ?)"
	_, err = utils.DB.Exec(query, loanId, payload.RejectionMessage, rejectionDate, payload.EmployeeID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Update Status
	query = "update loans set status = ? where id = ?"
	_, err = utils.DB.Exec(query, status, loanId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	resp := fiber.Map{
		"status":  "success",
		"message": "Loan successfully rejected",
		"data": fiber.Map{
			"loan_id":           loanId,
			"borrower_id":       loan.BorrowerID,
			"rejection_date":    payload.RejectionDate,
			"employee_id":       payload.EmployeeID,
			"rejection_message": payload.RejectionMessage,
			"loan_status":       status,
		},
	}

	return c.JSON(resp)
}
