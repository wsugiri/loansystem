package loans

import (
	"fmt"
	"log"

	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func CheckLoan(loanId int) (models.Loan, error) {
	var loan models.Loan
	var query = `
select a.id, a.borrower_id, a.principal_amount, a.rate, a.total_loan, a.instalment, a.duration_weeks, a.status, a.agreement_url
	 , sum(ifnull(b.amount, 0)) as invested_amount
  from loans a
  left join investments b on b.loan_id = a.id 
 where a.id = ?
 group by a.id`

	err := utils.DB.QueryRow(query, loanId).Scan(
		&loan.ID,
		&loan.BorrowerID,
		&loan.PrincipalAmount,
		&loan.Rate,
		&loan.TotalLoan,
		&loan.Instalment,
		&loan.DurationWeek,
		&loan.Status,
		&loan.AgreementUrl,
		&loan.InvestedAmount,
	)

	return loan, err
}

func GetPayments(loanId int) ([]models.Instalment, error) {
	var query = `select week, amount, due_date, is_paid, ifnull(payment_date, '') payment_date from payments where loan_id = ?`

	rows, err := utils.DB.Query(query, loanId)

	fmt.Println(loanId)

	var schedules []models.Instalment
	for rows.Next() {
		var sched models.Instalment
		if err := rows.Scan(&sched.Week, &sched.Amount, &sched.DueDate, &sched.IsPaid, &sched.PaymentDate); err != nil {
			log.Fatal(err)
		}
		schedules = append(schedules, sched)
	}

	return schedules, err
}
