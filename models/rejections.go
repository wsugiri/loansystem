package models

type Rejection struct {
	ID              int
	LoanID          int
	RejectionDate   string
	RejectionReason string
	RejectedBy      int
	Comment         string
}
