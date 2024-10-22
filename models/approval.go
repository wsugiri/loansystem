package models

type Approval struct {
	ID              int
	LoanID          int
	ApprovalDate    string
	PictureProofUrl string
	ApprovalBy      int
	Comment         string
}
