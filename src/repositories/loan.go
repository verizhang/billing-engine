package repositories

import (
	"context"
	"github.com/verizhang/billing-engine/src/entities"
	"gorm.io/gorm"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *entities.Loan) error
	GetActiveLoansByUserID(ctx context.Context, userID string) ([]*entities.Loan, error)
	UpdateIsActiveLoanByID(ctx context.Context, ID string, isActive bool) error
}

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{
		db: db,
	}
}

func (r *loanRepository) CreateLoan(ctx context.Context, loan *entities.Loan) error {
	if err := r.db.Create(loan).Error; err != nil {
		return err
	}

	return nil
}

func (r *loanRepository) GetActiveLoansByUserID(ctx context.Context, userID string) ([]*entities.Loan, error) {
	var loans []*entities.Loan
	err := r.db.Where("user_id = ? AND is_active = true", userID).Find(&loans).Error
	if err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *loanRepository) UpdateIsActiveLoanByID(ctx context.Context, ID string, isActive bool) error {
	if err := r.db.Model(&entities.Loan{}).Where("id = ?", ID).Update("is_active", isActive).Error; err != nil {
		return err
	}

	return nil
}
