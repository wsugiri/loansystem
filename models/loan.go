package models

type Loan struct {
	ID              int
	BorrowerID      int
	PrincipalAmount float32
	InvestedAmount  float32
	Rate            float32
	TotalLoan       float32
	Instalment      float32
	DurationWeek    int
	Status          string
	AgreementUrl    string
}

type LoanOutstanding struct {
	ID              int
	BorrowerID      int
	PrincipalAmount float32
	TotalLoan       float32
	OutstandingLoan float32
	DurationWeek    int
}
