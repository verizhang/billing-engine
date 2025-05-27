package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/verizhang/billing-engine/src/entities"
	mocks "github.com/verizhang/billing-engine/src/repositories/mocks"
	"github.com/verizhang/billing-engine/src/services"
	"github.com/verizhang/billing-engine/src/utils/errorhandler"
	"gorm.io/gorm"
)

func TestLoanService_CreateLoan(t *testing.T) {
	// Helper function to create service with mocks
	createService := func(uow *mocks.UnitOfWork, loanRepo *mocks.LoanRepository, paymentRepo *mocks.PaymentRepository) services.LoanService {
		if uow == nil {
			uow = new(mocks.UnitOfWork)
		}
		return services.NewLoanService(uow, loanRepo, paymentRepo)
	}

	t.Run("success create loan", func(t *testing.T) {
		// Setup mocks
		uow := new(mocks.UnitOfWork)
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)

		// Mock transaction
		mockTx := &gorm.DB{}

		// Mock expectations
		uow.On("Begin", mock.Anything).Return(mockTx, nil)
		uow.On("Commit", mockTx).Return(nil)
		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{}, nil)

		// Mock repository creation within UoW
		uow.On("LoanRepository", mockTx).Return(loanRepo)
		uow.On("PaymentRepository", mockTx).Return(paymentRepo)

		// Mock repository calls
		loanRepo.On("CreateLoan", mock.Anything, mock.AnythingOfType("*entities.Loan")).Return(nil)
		paymentRepo.On("CreatePayments", mock.Anything, mock.AnythingOfType("[]*entities.Payment")).Return(nil)

		// Execute
		service := createService(uow, loanRepo, paymentRepo)
		err := service.CreateLoan(context.Background(), "user1")

		// Assert
		assert.NoError(t, err)
		uow.AssertExpectations(t)
		loanRepo.AssertExpectations(t)
		paymentRepo.AssertExpectations(t)
	})

	t.Run("error when begin transaction fails", func(t *testing.T) {
		uow := new(mocks.UnitOfWork)
		loanRepo := new(mocks.LoanRepository)

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{}, nil)
		uow.On("Begin", mock.Anything).Return(nil, errors.New("transaction error"))

		service := createService(uow, loanRepo, nil)
		err := service.CreateLoan(context.Background(), "user1")

		assert.Error(t, err)
		assert.Equal(t, errorhandler.InternalServerError, errors.Unwrap(err))
	})

	t.Run("error when create loan fails - should rollback", func(t *testing.T) {
		uow := new(mocks.UnitOfWork)
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)
		mockTx := &gorm.DB{}

		uow.On("Begin", mock.Anything).Return(mockTx, nil)
		uow.On("Rollback", mockTx).Return(nil)
		uow.On("LoanRepository", mockTx).Return(loanRepo)
		// Mock repository creation within UoW
		uow.On("LoanRepository", mockTx).Return(loanRepo)
		uow.On("PaymentRepository", mockTx).Return(paymentRepo)

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{}, nil)
		loanRepo.On("CreateLoan", mock.Anything, mock.Anything).Return(errors.New("create error"))

		service := createService(uow, loanRepo, paymentRepo)
		err := service.CreateLoan(context.Background(), "user1")

		assert.Error(t, err)
		uow.AssertCalled(t, "Rollback", mockTx)
	})

	t.Run("error when create payments fails - should rollback", func(t *testing.T) {
		uow := new(mocks.UnitOfWork)
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)
		mockTx := &gorm.DB{}

		uow.On("Begin", mock.Anything).Return(mockTx, nil)
		uow.On("Rollback", mockTx).Return(nil)
		uow.On("LoanRepository", mockTx).Return(loanRepo)
		uow.On("PaymentRepository", mockTx).Return(paymentRepo)
		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{}, nil)
		loanRepo.On("CreateLoan", mock.Anything, mock.Anything).Return(nil)
		paymentRepo.On("CreatePayments", mock.Anything, mock.Anything).Return(errors.New("payment error"))

		service := createService(uow, loanRepo, paymentRepo)
		err := service.CreateLoan(context.Background(), "user1")

		assert.Error(t, err)
		uow.AssertCalled(t, "Rollback", mockTx)
	})
}

func TestLoanService_GetOutstanding(t *testing.T) {
	createService := func(loanRepo *mocks.LoanRepository, paymentRepo *mocks.PaymentRepository) services.LoanService {
		return services.NewLoanService(nil, loanRepo, paymentRepo)
	}

	t.Run("success with no payments", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)

		loan := &entities.Loan{
			ID:       "loan1",
			Amount:   1000,
			Interest: 100,
		}

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return([]*entities.Payment{}, nil)

		service := createService(loanRepo, paymentRepo)
		result, err := service.GetOutstanding(context.Background(), "user1")

		assert.NoError(t, err)
		assert.Equal(t, float64(1100), result.Outstanding)
	})

	t.Run("error when get payments fails", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)

		loan := &entities.Loan{ID: "loan1"}
		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(nil, errors.New("db error"))

		service := createService(loanRepo, paymentRepo)
		_, err := service.GetOutstanding(context.Background(), "user1")

		assert.Error(t, err)
		assert.Equal(t, errorhandler.InternalServerError, errors.Unwrap(err))
	})
}

func TestLoanService_IsDelinquent(t *testing.T) {
	createService := func(loanRepo *mocks.LoanRepository, paymentRepo *mocks.PaymentRepository) services.LoanService {
		return services.NewLoanService(nil, loanRepo, paymentRepo)
	}

	t.Run("delinquent when payment overdue", func(t *testing.T) {
		loanRepo := new(mocks.LoanRepository)
		paymentRepo := new(mocks.PaymentRepository)

		loan := &entities.Loan{ID: "loan1"}
		dueDate := time.Now().AddDate(0, 0, -16) // 15 days ago (overdue)
		payments := []*entities.Payment{
			{PaidAt: nil, EndAt: &dueDate},
		}

		loanRepo.On("GetActiveLoansByUserID", mock.Anything, "user1").Return([]*entities.Loan{loan}, nil)
		paymentRepo.On("GetPaymentByLoanID", mock.Anything, "loan1").Return(payments, nil)

		service := createService(loanRepo, paymentRepo)
		result, err := service.IsDelinquent(context.Background(), "user1")

		assert.NoError(t, err)
		assert.True(t, result.IsDelinquent)
	})
}
