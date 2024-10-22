package models

type Disbursement struct {
	ID                 int
	LoanID             int
	FieldOfficerID     int
	DisbursementDate   string
	AgreementSignedUrl string
}
