package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/verizhang/billing-engine/src/entities"
	"github.com/verizhang/billing-engine/src/utils/errorhandler"
	"time"
)

func (s *loanService) generatePayments(loanID string, now time.Time) []*entities.Payment {
	var payments []*entities.Payment
	amount := entities.LOAN_AMOUNT/entities.PAYMENT_WEEKS + (entities.LOAN_AMOUNT*entities.LOAN_INTEREST_RATE)/entities.PAYMENT_WEEKS
	for i := 0; i < entities.PAYMENT_WEEKS; i++ {
		startAt := now.AddDate(0, 0, 7*i)
		endAt := startAt.AddDate(0, 0, 7).Add(-time.Nanosecond)
		uuid, _ := uuid.NewUUID()
		payments = append(payments, &entities.Payment{
			ID:      uuid.String(),
			LoanID:  loanID,
			Amount:  amount,
			StartAt: &startAt,
			EndAt:   &endAt,
			PaidAt:  nil,
		})
	}

	return payments
}

func (s *loanService) getActiveLoan(ctx context.Context, userID string) (*entities.Loan, error) {
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

func (s *loanService) compareDelinquent(payments []*entities.Payment) bool {
	now := time.Now()

	lastPaidIdx := 0
	for i, payment := range payments {
		if payment.PaidAt != nil {
			lastPaidIdx = i
			break
		}
	}

	lastUnpaidIdx := lastPaidIdx
	if len(payments)-1 > lastPaidIdx+1 {
		lastUnpaidIdx = lastPaidIdx + 1
	}

	lastPaid := payments[lastUnpaidIdx]

	dueDate := lastPaid.EndAt.AddDate(0, 0, 14)
	return now.After(dueDate)
}
