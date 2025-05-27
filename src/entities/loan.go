package entities

import "time"

const (
	LOAN_AMOUNT        = 5000000
	LOAN_INTEREST_RATE = 0.10
)

type Loan struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Amount    float64    `json:"amount"`
	Interest  float64    `json:"interest"`
	IsActive  bool       `json:"is_active"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy string     `json:"updated_by"`
	DeletedBy string     `json:"deleted_by"`
}

type Outstanding struct {
	Outstanding float64
}

type IsDelinquent struct {
	IsDelinquent bool
}
