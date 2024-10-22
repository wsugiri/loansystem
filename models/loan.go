package models

type Loan struct {
	ID              int
	BorrowerID      int
	PrincipalAmount float32
	Rate            float32
	TotalLoan       float32
	Instalment      float32
	Status          string
}
