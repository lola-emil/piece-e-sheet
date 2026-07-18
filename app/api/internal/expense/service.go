package expense

import (
	"context"
	"database/sql"
	"errors"
)

type service struct {
	expenseRepo ExpenseRepository
}

type ExpenseService interface {
	FindAll(ctx context.Context, userID string, filter ExpenseFilter) ([]Expense, error)
	FindByID(ctx context.Context, id string, userID string) (*Expense, error)
	Create(ctx context.Context, userID string, req *CreateExpenseRequest) (*Expense, error)
	Update(ctx context.Context, id string, userID string, req *UpdateExpenseRequest) (*Expense, error)
	Delete(ctx context.Context, id string, userID string) error
}

func NewExpenseService(expenseRepo ExpenseRepository) ExpenseService {
	return &service{expenseRepo: expenseRepo}
}

func (s *service) FindAll(ctx context.Context, userID string, filter ExpenseFilter) ([]Expense, error) {
	return s.expenseRepo.FindAll(ctx, userID, filter)
}

func (s *service) FindByID(ctx context.Context, id string, userID string) (*Expense, error) {
	expense, err := s.expenseRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if expense == nil {
		return nil, errors.New("expense not found")
	}
	return expense, nil
}

func (s *service) Create(ctx context.Context, userID string, req *CreateExpenseRequest) (*Expense, error) {
	if req.Description == "" {
		return nil, errors.New("description is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	expense := &Expense{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Amount:      req.Amount,
		OccurredAt:  req.OccurredAt,
	}

	err := s.expenseRepo.Insert(ctx, expense)
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *service) Update(ctx context.Context, id string, userID string, req *UpdateExpenseRequest) (*Expense, error) {
	if req.Description == "" {
		return nil, errors.New("description is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	expense, err := s.expenseRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if expense == nil {
		return nil, errors.New("expense not found")
	}

	expense.CategoryID = req.CategoryID
	expense.Description = req.Description
	expense.Amount = req.Amount
	expense.OccurredAt = req.OccurredAt

	err = s.expenseRepo.UpdateByID(ctx, expense)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("expense not found or conflict")
		}
		return nil, err
	}
	return expense, nil
}

func (s *service) Delete(ctx context.Context, id string, userID string) error {
	err := s.expenseRepo.DeleteByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("expense not found")
		}
		return err
	}
	return nil
}
