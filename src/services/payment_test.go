package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/verizhang/billing-engine/config"
	"github.com/verizhang/billing-engine/src/entities"
	mocks "github.com/verizhang/billing-engine/src/repositories/mocks"
	"github.com/verizhang/billing-engine/src/services"
	"github.com/verizhang/billing-engine/src/utils/errorhandler"
	"gorm.io/gorm"
)

func TestPaymentService_MakePayment(t *testing.T) {
	// Helper function to create service with mocks
	createService := func(
		cfg config.Config,
		uow *mocks.UnitOfWork,
		paymentRepo *mocks.PaymentRepository,
		loanRepo *mocks.LoanRepository,
	) services.PaymentService {
		return services.NewPaymentService(cfg, paymentRepo, loanRepo, uow)
	}

	t.Run("success make payment - not last payment", func(t *testing.T) {
		// Setup mocks
		uow := new(mocks.UnitOfWork)
		paymentRepo := new(mocks.PaymentRepository)
		loanRepo := new(mocks.LoanRepository)
		mockTx := &gorm.DB{}

		// Test data
		now := time.Now()
		lastDay := now.AddDate(0, 0, -1)
		nextEndAt := now.AddDate(0, 0, 7)
		payment2EndAt := nextEndAt.AddDate(0, 0, 7)
		loan := &entities.Loan{ID: "loan1", UserID: "user1", IsActive: true}
		payments := []*entities.Payment{
			{
				ID:      "payment1",
				LoanID:  "loan1",
				PaidAt:  nil,
				StartAt: &lastDay,
				EndAt:   &nextEndAt,
			},
			{
				ID:      "payment2",
				LoanID:  "loan1",
				PaidAt:  nil,
				StartAt: &nextEndAt,
				EndAt:   &payment2EndAt,
			},
		}

		// Mock expectations
		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(payments, nil)
		paymentRepo.On("UpdatePaidAtPayment", mock.Anything, "payment1", mock.Anything).Return(nil)
		uow.On("Begin", mock.Anything).Return(mockTx, nil)
		uow.On("Commit", mockTx).Return(nil)
		uow.On("PaymentRepository", mockTx).Return(paymentRepo)
		uow.On("LoanRepository", mockTx).Return(loanRepo)

		// Execute
		service := createService(config.Config{}, uow, paymentRepo, loanRepo)
		err := service.MakePayment(context.Background(), "user1")

		// Assert
		assert.NoError(t, err)
		uow.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		paymentRepo.AssertExpectations(t)
	})

	t.Run("success make payment - last payment", func(t *testing.T) {
		// Setup mocks
		uow := new(mocks.UnitOfWork)
		paymentRepo := new(mocks.PaymentRepository)
		loanRepo := new(mocks.LoanRepository)
		mockTx := &gorm.DB{}

		// Test data
		now := time.Now()
		loan := &entities.Loan{ID: "loan1", UserID: "user1", IsActive: true}
		payments := []*entities.Payment{
			{
				ID:      "payment1",
				LoanID:  "loan1",
				PaidAt:  &now,
				StartAt: &now,
				EndAt:   &now,
			},
			{
				ID:      "payment2",
				LoanID:  "loan1",
				PaidAt:  nil,
				StartAt: &now,
				EndAt:   &now,
			},
		}

		// Mock expectations
		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(payments, nil)
		uow.On("Begin", mock.Anything).Return(mockTx, nil)
		uow.On("Commit", mockTx).Return(nil)
		uow.On("PaymentRepository", mockTx).Return(paymentRepo)
		uow.On("LoanRepository", mockTx).Return(loanRepo)
		paymentRepo.On("UpdatePaidAtPayment", mock.Anything, "payment2", mock.Anything).Return(nil)
		loanRepo.On("UpdateIsActiveLoanByID", mock.Anything, "loan1", false).Return(nil)

		// Execute
		service := createService(config.Config{}, uow, paymentRepo, loanRepo)
		err := service.MakePayment(context.Background(), "user1")

		// Assert
		assert.NoError(t, err)
		uow.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		paymentRepo.AssertExpectations(t)
	})

	t.Run("error - no active loan found", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{}, nil)

		service := createService(config.Config{}, nil, nil, loanRepo)
		err := service.MakePayment(context.Background(), "user1")

		assert.Error(t, err)
		assert.Equal(t, errorhandler.BadRequestError, errors.Unwrap(err))
	})

	t.Run("error - all payments already paid", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)

		now := time.Now()
		loan := &entities.Loan{ID: "loan1", UserID: "user1", IsActive: true}
		payments := []*entities.Payment{
			{
				ID:      "payment1",
				LoanID:  "loan1",
				PaidAt:  &now,
				StartAt: &now,
				EndAt:   &now,
			},
		}

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(payments, nil)

		service := createService(config.Config{}, nil, paymentRepo, loanRepo)
		err := service.MakePayment(context.Background(), "user1")

		assert.Error(t, err)
		assert.Equal(t, errorhandler.BadRequestError, errors.Unwrap(err))
	})

	t.Run("error - begin transaction fails", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)
		uow := new(mocks.UnitOfWork)

		now := time.Now()
		loan := &entities.Loan{ID: "loan1", UserID: "user1", IsActive: true}
		payments := []*entities.Payment{
			{
				ID:      "payment1",
				LoanID:  "loan1",
				PaidAt:  nil,
				StartAt: &now,
				EndAt:   &now,
			},
		}

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(payments, nil)
		uow.On("Begin", mock.Anything).Return(nil, errors.New("transaction error"))

		service := createService(config.Config{}, uow, paymentRepo, loanRepo)
		err := service.MakePayment(context.Background(), "user1")

		assert.Error(t, err)
		assert.Equal(t, errorhandler.InternalServerError, errors.Unwrap(err))
	})

	t.Run("error - update payment fails", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)
		uow := new(mocks.UnitOfWork)
		mockTx := &gorm.DB{}

		now := time.Now()
		loan := &entities.Loan{ID: "loan1", UserID: "user1", IsActive: true}
		payments := []*entities.Payment{
			{
				ID:      "payment1",
				LoanID:  "loan1",
				PaidAt:  nil,
				StartAt: &now,
				EndAt:   &now,
			},
		}

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(payments, nil)
		paymentRepo.On("UpdatePaidAtPayment", mock.Anything, "payment1", mock.Anything).Return(errors.New("update error"))
		uow.On("Begin", mock.Anything).Return(mockTx, nil)
		uow.On("Rollback", mockTx).Return(nil)
		uow.On("PaymentRepository", mockTx).Return(paymentRepo)
		uow.On("LoanRepository", mockTx).Return(loanRepo)

		service := createService(config.Config{}, uow, paymentRepo, loanRepo)
		err := service.MakePayment(context.Background(), "user1")

		assert.Error(t, err)
		assert.Equal(t, errorhandler.InternalServerError, errors.Unwrap(err))
		uow.AssertCalled(t, "Rollback", mockTx)
	})

	t.Run("error - update loan status fails (last payment)", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)
		uow := new(mocks.UnitOfWork)
		mockTx := &gorm.DB{}

		now := time.Now()
		loan := &entities.Loan{ID: "loan1", UserID: "user1", IsActive: true}
		payments := []*entities.Payment{
			{
				ID:      "payment1",
				LoanID:  "loan1",
				PaidAt:  nil,
				StartAt: &now,
				EndAt:   &now,
			},
		}

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(payments, nil)
		uow.On("Begin", mock.Anything).Return(mockTx, nil)
		uow.On("Rollback", mockTx).Return(nil)
		uow.On("PaymentRepository", mockTx).Return(paymentRepo)
		uow.On("LoanRepository", mockTx).Return(loanRepo)
		paymentRepo.On("UpdatePaidAtPayment", mock.Anything, "payment1", mock.Anything).Return(nil)
		loanRepo.On("UpdateIsActiveLoanByID", mock.Anything, "loan1", false).Return(errors.New("update error"))

		service := createService(config.Config{}, uow, paymentRepo, loanRepo)
		err := service.MakePayment(context.Background(), "user1")

		assert.Error(t, err)
		assert.Equal(t, errorhandler.InternalServerError, errors.Unwrap(err))
		uow.AssertCalled(t, "Rollback", mockTx)
	})

	t.Run("error - commit transaction fails", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)
		uow := new(mocks.UnitOfWork)
		mockTx := &gorm.DB{}

		now := time.Now()
		loan := &entities.Loan{ID: "loan1", UserID: "user1", IsActive: true}
		payments := []*entities.Payment{
			{
				ID:      "payment1",
				LoanID:  "loan1",
				PaidAt:  nil,
				StartAt: &now,
				EndAt:   &now,
			},
		}

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		loanRepo.On("UpdateIsActiveLoanByID", mock.Anything, "loan1", false).Return(nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(payments, nil)
		paymentRepo.On("UpdatePaidAtPayment", mock.Anything, "payment1", mock.Anything).Return(nil)
		uow.On("Begin", mock.Anything).Return(mockTx, nil)
		uow.On("Commit", mockTx).Return(errors.New("commit error"))
		uow.On("PaymentRepository", mockTx).Return(paymentRepo)
		uow.On("LoanRepository", mockTx).Return(loanRepo)

		service := createService(config.Config{}, uow, paymentRepo, loanRepo)
		err := service.MakePayment(context.Background(), "user1")

		assert.Error(t, err)
		assert.Equal(t, errorhandler.InternalServerError, errors.Unwrap(err))
	})
}
