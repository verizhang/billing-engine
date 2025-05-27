package handlers

import (
	"context"
	loanpb "github.com/verizhang/billing-engine/contracts/pb/loan"
	"github.com/verizhang/billing-engine/src/services"
	"github.com/verizhang/billing-engine/src/utils/errorhandler"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LoanHandler struct {
	loanpb.UnimplementedLoanServer
	svc services.LoanService
}

func NewLoanHandler(svc services.LoanService) *LoanHandler {
	return &LoanHandler{
		svc: svc,
	}
}

func (h *LoanHandler) CreateLoan(ctx context.Context, req *loanpb.CreateLoanRequest) (*emptypb.Empty, error) {
	err := h.svc.CreateLoan(ctx, req.UserId)
	if err != nil {
		return nil, errorhandler.TranslateTogRPCError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *LoanHandler) GetOutstanding(ctx context.Context, req *loanpb.GetOutstandingRequest) (*loanpb.GetOutstandingResponse, error) {
	resp, err := h.svc.GetOutstanding(ctx, req.UserId)
	if err != nil {
		return nil, errorhandler.TranslateTogRPCError(err)
	}

	return &loanpb.GetOutstandingResponse{Outstanding: float32(resp.Outstanding)}, nil
}

func (h *LoanHandler) IsDelinquent(ctx context.Context, req *loanpb.GetIsDelinquentRequest) (*loanpb.GetIsDelinquentResponse, error) {
	resp, err := h.svc.IsDelinquent(ctx, req.UserId)
	if err != nil {
		return nil, errorhandler.TranslateTogRPCError(err)
	}

	return &loanpb.GetIsDelinquentResponse{IsDelinquent: resp.IsDelinquent}, nil
}
