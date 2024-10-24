package loans

import (
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
