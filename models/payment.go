package models

type Payment struct {
	ID          int
	LoanID      int
	Week        int
	Amount      float32
	DueDate     string
	PaymentDate string
}

type Instalment struct {
	Week    int
	Amount  float32
	DueDate string
}
