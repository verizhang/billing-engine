package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/verizhang/billing-engine/config"
	"github.com/verizhang/billing-engine/src/repositories"
	"github.com/verizhang/billing-engine/src/utils/errorhandler"
	"time"
)

type PaymentService interface {
	MakePayment(ctx context.Context, userID string) error
}

type paymentService struct {
	cfg         config.Config
	paymentRepo repositories.PaymentRepository
	uow         repositories.UnitOfWork
	loanRepo    repositories.LoanRepository
}

func NewPaymentService(cfg config.Config, paymentRepo repositories.PaymentRepository, loanRepo repositories.LoanRepository, uow repositories.UnitOfWork) PaymentService {
	return &paymentService{
		cfg:         cfg,
		paymentRepo: paymentRepo,
		loanRepo:    loanRepo,
		uow:         uow,
	}
}

func (s *paymentService) MakePayment(ctx context.Context, userID string) error {
	loan, err := s.getActiveLoan(ctx, userID)
	if err != nil {
		return err
	}
	now := time.Now()

	payments, err := s.paymentRepo.GetPaymentByLoanID(ctx, loan.ID)
	unpaidPayment := s.getEligiblePayment(payments, &now)
	if unpaidPayment == nil {
		return fmt.Errorf("%w: %s", errorhandler.BadRequestError, errors.New("all loans have been paid off").Error())
	}
	isLastPayment := s.isLastPayment(payments, unpaidPayment.ID)

	tx, err := s.uow.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	loanRepo := s.uow.LoanRepository(tx)
	paymentRepo := s.uow.PaymentRepository(tx)

	err = paymentRepo.UpdatePaidAtPayment(ctx, unpaidPayment.ID, &now)
	if err != nil {
		s.uow.Rollback(tx)
		return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	if isLastPayment {
		err = loanRepo.UpdateIsActiveLoanByID(ctx, loan.ID, false)
		if err != nil {
			s.uow.Rollback(tx)
			return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
		}
	}

	err = s.uow.Commit(tx)
	if err != nil {
		return fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	return nil
}
