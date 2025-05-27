package repositories

import (
	"context"
	"github.com/verizhang/billing-engine/src/entities"
	"gorm.io/gorm"
	"time"
)

type PaymentRepository interface {
	CreatePayments(ctx context.Context, payments []*entities.Payment) error
	UpdatePaidAtPayment(ctx context.Context, ID string, paidAt *time.Time) error
	GetPaymentByLoanID(ctx context.Context, loanID string) ([]*entities.Payment, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) CreatePayments(ctx context.Context, payments []*entities.Payment) error {
	if err := r.db.Create(payments).Error; err != nil {
		return err
	}
	return nil
}

func (r *paymentRepository) UpdatePaidAtPayment(ctx context.Context, ID string, paidAt *time.Time) error {
	if err := r.db.Model(&entities.Payment{}).Where("id = ?", ID).Update("paid_at", paidAt).Error; err != nil {
		return err
	}
	return nil
}

func (r *paymentRepository) GetPaymentByLoanID(ctx context.Context, loanID string) ([]*entities.Payment, error) {
	var payments []*entities.Payment
	if err := r.db.Where("loan_id = ?", loanID).Order("start_at ASC").Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}
