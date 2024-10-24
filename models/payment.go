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
	Week        int     `json:"week"`
	Amount      float32 `json:"amount"`
	DueDate     string  `json:"due_date"`
	PaymentDate string  `json:"payment_date"`
	IsPaid      bool    `json:"paid"`
}
