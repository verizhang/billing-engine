package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/verizhang/billing-engine/src/entities"
	"github.com/verizhang/billing-engine/src/utils/errorhandler"
	"time"
)

func (s *paymentService) getActiveLoan(ctx context.Context, userID string) (*entities.Loan, error) {
	loans, err := s.loanRepo.GetActiveLoansByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errorhandler.InternalServerError, err.Error())
	}

	var loan *entities.Loan
	if len(loans) > 0 {
		loan = loans[0]
	} else {
		return nil, fmt.Errorf("%w: %s", errorhandler.BadRequestError, errors.New("Active loan not found"))
	}

	return loan, nil
}

func (s *paymentService) getEligiblePayment(payments []*entities.Payment, now *time.Time) *entities.Payment {
	for _, payment := range payments {

		if now.After(*payment.StartAt) && payment.PaidAt == nil {
			return payment
		}
	}

	return nil
}

func (s *paymentService) isLastPayment(payments []*entities.Payment, ID string) bool {
	lastPayment := payments[len(payments)-1]
	if lastPayment.ID == ID {
		return true
	}
	return false
}
