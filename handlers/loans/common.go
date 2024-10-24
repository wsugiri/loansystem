package loans

import (
	"github.com/wsugiri/loansystem/models"
	"github.com/wsugiri/loansystem/utils"
)

func CheckLoan(loanId int) (models.Loan, error) {
	var loan models.Loan
	var query = `
select a.id, a.borrower_id, a.principal_amount, a.status, a.rate
	 , sum(ifnull(b.amount, 0)) as invested_amount
  from loans a
  left join investments b on b.loan_id = a.id 
 where a.id = ?
 group by a.id`

	err := utils.DB.QueryRow(query, loanId).Scan(&loan.ID, &loan.BorrowerID, &loan.PrincipalAmount, &loan.Status, &loan.Rate, &loan.InvestedAmount)

	return loan, err
}
