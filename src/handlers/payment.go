package handlers

import (
	"context"
	paymentpb "github.com/verizhang/billing-engine/contracts/pb/payment"
	"github.com/verizhang/billing-engine/src/services"
	"github.com/verizhang/billing-engine/src/utils/errorhandler"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServer
	svc services.PaymentService
}

func NewPaymentHandler(svc services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		svc: svc,
	}
}

func (h *PaymentHandler) MakePayment(ctx context.Context, req *paymentpb.MakePaymentRequest) (*emptypb.Empty, error) {
	err := h.svc.MakePayment(ctx, req.UserId)
	if err != nil {
		return nil, errorhandler.TranslateTogRPCError(err)
	}

	return &emptypb.Empty{}, nil
}
