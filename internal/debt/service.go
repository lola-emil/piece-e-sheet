package debt

import (
	"api/internal/store"
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type service struct{ repo DebtRepository }

type DebtService interface {
	FindAll(ctx context.Context, userID string) ([]Debt, error)
	Create(ctx context.Context, userID string, req *CreateDebtRequest) (*Debt, error)
	Update(ctx context.Context, id, userID string, req *UpdateDebtRequest) (*Debt, error)
	Delete(ctx context.Context, id, userID string) error
}

func NewDebtService(repo DebtRepository) DebtService { return &service{repo} }

func (s *service) FindAll(ctx context.Context, userID string) ([]Debt, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *service) Create(ctx context.Context, userID string, req *CreateDebtRequest) (*Debt, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request body is required", store.ErrInvalidInput)
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	d := &Debt{
		UserID:      userID,
		PersonName:  req.PersonName,
		Description: req.Description,
		Amount:      req.Amount,
		Type:        req.Type,
		DueDate:     req.DueDate,
	}

	if err := s.repo.Insert(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *service) Update(ctx context.Context, id, userID string, req *UpdateDebtRequest) (*Debt, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: request body is required", store.ErrInvalidInput)
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	d, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	d.PersonName = req.PersonName
	d.Description = req.Description
	d.Amount = req.Amount
	d.RemainingAmount = req.RemainingAmount
	d.Type = req.Type
	d.Status = req.Status
	d.DueDate = req.DueDate

	d.Revision = req.Revision

	if req.Status == "paid" {
		d.RemainingAmount = decimal.Zero
	}

	normalizeSettlement(d)

	if err := s.repo.UpdateByID(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *service) Delete(ctx context.Context, id, userID string) error {
	return s.repo.DeleteByID(ctx, id, userID)
}

func normalizeSettlement(d *Debt) {
	switch {
	case d.RemainingAmount.IsZero():
		d.Status = "paid"
		if d.SettledAt == nil {
			now := time.Now()
			d.SettledAt = &now
		}

	case d.RemainingAmount.Equal(d.Amount):
		d.Status = "pending"
		d.SettledAt = nil

	default:
		d.Status = "partial"
		d.SettledAt = nil
	}
}
