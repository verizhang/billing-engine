package repositories

import (
	"context"
	"gorm.io/gorm"
)

type UnitOfWork interface {
	Begin(ctx context.Context) (*gorm.DB, error)
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error
	LoanRepository(tx *gorm.DB) LoanRepository
	PaymentRepository(tx *gorm.DB) PaymentRepository
}

type unitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

func (u *unitOfWork) Begin(ctx context.Context) (*gorm.DB, error) {
	tx := u.db.Begin().WithContext(ctx)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return tx, nil
}

func (u *unitOfWork) Commit(tx *gorm.DB) error {
	err := tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (u *unitOfWork) Rollback(tx *gorm.DB) error {
	err := tx.Rollback().Error
	if err != nil {
		return err
	}
	return nil
}

func (u *unitOfWork) LoanRepository(tx *gorm.DB) LoanRepository {
	return NewLoanRepository(tx)
}

func (u *unitOfWork) PaymentRepository(tx *gorm.DB) PaymentRepository {
	return NewPaymentRepository(tx)
}
