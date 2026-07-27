package payment

import (
	"context"
	"errors"
)

type service struct {
	repo PaymentRepository
}

type PaymentService interface {
	GetByDebtID(ctx context.Context, debtID, userID string) ([]Payment, error)
	Record(ctx context.Context, userID string, req *CreatePaymentRequest) (*Payment, error)
}

func NewPaymentService(repo PaymentRepository) PaymentService { return &service{repo} }

func (s *service) GetByDebtID(ctx context.Context, debtID, userID string) ([]Payment, error) {
	return s.repo.FindByDebtID(ctx, debtID, userID)
}

func (s *service) Record(ctx context.Context, userID string, req *CreatePaymentRequest) (*Payment, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	p := &Payment{
		UserID: userID,
		DebtID: req.DebtID,
		Amount: req.Amount,
		Note:   req.Note,
	}

	err := s.repo.RecordPayment(ctx, p)
	return p, err
}
