package expense

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	expenseRepo ExpenseRepository
}

type ExpenseService interface {
	FindAll(ctx context.Context, userID uuid.UUID) ([]Expense, error)
	FindByID(ctx context.Context, expenseID string) (*Expense, error)

	Insert(ctx context.Context, e *Expense) error
	UpdateByID(ctx context.Context, e *Expense) error
	DeleteByID(ctx context.Context, expenseID string) error
}

func NewExpenseService(expenseRepo ExpenseRepository) ExpenseService {
	return &service{expenseRepo: expenseRepo}
}

func (s *service) FindAll(ctx context.Context, userID uuid.UUID) ([]Expense, error) {
	return s.expenseRepo.FindAll(ctx, userID)
}

func (s *service) FindByID(ctx context.Context, expenseID string) (*Expense, error) {
	return s.expenseRepo.FindByID(ctx, expenseID)
}

func (s *service) Insert(ctx context.Context, e *Expense) error {
	return s.expenseRepo.Insert(ctx, e)
}

func (s *service) UpdateByID(ctx context.Context, e *Expense) error {
	return s.expenseRepo.UpdateByID(ctx, e)
}

func (s *service) DeleteByID(ctx context.Context, expenseID string) error {
	return s.expenseRepo.DeleteByID(ctx, expenseID)
}
