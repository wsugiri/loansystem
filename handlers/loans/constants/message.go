package constants

const (
	ErrEmployeeInvalid      = "The provided employee_id is invalid or does not exist. Please check and try again"
	ErrEmployeeNoAuthorized = "The provided employee ID is not authorized to approve this loan"

	ErrBorrowerInvalid = "The provided borrower_id is invalid or does not exist. Please check and try again"

	ErrLoanInvalid       = "The provided loan_id is invalid or does not exist. Please check and try again"
	ErrLoanNotInProposed = "Loan is not in an proposed state"
	ErrLoanApproved      = "Loan status already approved"
	ErrLoanRejected      = "Loan status already rejected"
)
