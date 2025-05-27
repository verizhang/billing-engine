package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/verizhang/billing-engine/src/entities"
	"github.com/verizhang/billing-engine/src/repositories"
	"github.com/verizhang/billing-engine/src/utils/errorhandler"
	"time"
)

type LoanService interface {
	CreateLoan(ctx context.Context, userID string) error
	GetOutstanding(ctx context.Context, userID string) (*entities.Outstanding, error)
	IsDelinquent(ctx context.Context, userID string) (*entities.IsDelinquent, error)
}

type loanService struct {
	uow         repositories.UnitOfWork
	loanRepo    repositories.LoanRepository
	paymentRepo repositories.PaymentRepository
}

func NewLoanService(uow repositories.UnitOfWork, loanRepo repositories.LoanRepository, paymentRepo repositories.PaymentRepository) LoanService {
	return &loanService{
		uow:         uow,
		loanRepo:    loanRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *loanService) CreateLoan(ctx context.Context, userID string) error {
	loans, err := s.loanRepo.GetActiveLoansByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	if len(loans) > 0 {
		return fmt.Errorf("%w: you already apply for loan", errorhandler.BadRequestError)
	}

	tx, err := s.uow.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	loanRepo := s.uow.LoanRepository(tx)
	paymentRepo := s.uow.PaymentRepository(tx)

	now := time.Now()
	loanID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	err = loanRepo.CreateLoan(ctx, &entities.Loan{
		ID:        loanID.String(),
		UserID:    userID,
		Amount:    entities.LOAN_AMOUNT,
		Interest:  entities.LOAN_AMOUNT * entities.LOAN_INTEREST_RATE,
		IsActive:  true,
		CreatedAt: &now,
	})
	if err != nil {
		s.uow.Rollback(tx)
		return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	err = paymentRepo.CreatePayments(ctx, s.generatePayments(loanID.String(), now))
	if err != nil {
		s.uow.Rollback(tx)
		return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	s.uow.Commit(tx)
	return nil
}

func (s *loanService) GetOutstanding(ctx context.Context, userID string) (*entities.Outstanding, error) {
	loan, err := s.getActiveLoan(ctx, userID)
	if err != nil {
		return nil, err
	}

	outstanding := loan.Amount + loan.Interest

	payments, err := s.paymentRepo.GetPaymentByLoanID(ctx, loan.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	for _, payment := range payments {
		if payment.PaidAt != nil {
			outstanding = outstanding - payment.Amount
		}
	}

	return &entities.Outstanding{
		Outstanding: outstanding,
	}, nil
}

func (s *loanService) IsDelinquent(ctx context.Context, userID string) (*entities.IsDelinquent, error) {
	loan, err := s.getActiveLoan(ctx, userID)
	if err != nil {
		return nil, err
	}

	payments, err := s.paymentRepo.GetPaymentByLoanID(ctx, loan.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	isDelinquent := s.compareDelinquent(payments)

	return &entities.IsDelinquent{IsDelinquent: isDelinquent}, nil
}
