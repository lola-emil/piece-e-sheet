// debt/service.go
package debt

import (
	"context"
	"database/sql"
	"errors"
	"time"
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
	if req.Amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}
	d := &Debt{UserID: userID, PersonName: req.PersonName, Description: req.Description, Amount: req.Amount, Type: req.Type, DueDate: req.DueDate}
	err := s.repo.Insert(ctx, d)
	return d, err
}

func (s *service) Update(ctx context.Context, id, userID string, req *UpdateDebtRequest) (*Debt, error) {
	d, err := s.repo.FindByID(ctx, id, userID)
	if err != nil || d == nil {
		return nil, errors.New("not found")
	}

	d.PersonName = req.PersonName
	d.Description = req.Description
	d.Amount = req.Amount
	d.RemainingAmount = req.RemainingAmount
	d.Type = req.Type
	d.Status = req.Status
	d.DueDate = req.DueDate

	if req.Status == "paid" {
		now := time.Now()
		d.SettledAt = &now
		d.RemainingAmount = 0
	} else if req.Status != "paid" {
		d.SettledAt = nil
	}

	err = s.repo.UpdateByID(ctx, d)
	return d, err
}

func (s *service) Delete(ctx context.Context, id, userID string) error {
	err := s.repo.DeleteByID(ctx, id, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("not found")
	}
	return err
}
