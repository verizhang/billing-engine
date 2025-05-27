package entities

import "time"

const PAYMENT_WEEKS = 50

type Payment struct {
	ID        string     `json:"id"`
	LoanID    string     `json:"loan_id"`
	Amount    float64    `json:"amount"`
	StartAt   *time.Time `json:"start_date"`
	EndAt     *time.Time `json:"end_date"`
	PaidAt    *time.Time `json:"paid_at"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedBy *int64     `json:"created_by"`
	UpdatedBy *int64     `json:"updated_by"`
	DeletedBy *int64     `json:"deleted_by"`
}
