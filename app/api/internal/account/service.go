package account

import (
	"context"
	"database/sql"
	"errors"
)

type service struct {
	accountRepo AccountRepository
}

type AccountService interface {
	FindAll(ctx context.Context, userID string) ([]Account, error)
	FindByID(ctx context.Context, id, userID string) (*Account, error)
	Create(ctx context.Context, userID string, req *CreateAccountRequest) (*Account, error)
	Update(ctx context.Context, id string, userID string, req *UpdateAccountRequest) (*Account, error)
	Delete(ctx context.Context, id, userId string) error
}

func NewAccountService(accountRepo AccountRepository) AccountService {
	return &service{accountRepo: accountRepo}
}

func (s *service) FindAll(ctx context.Context, userID string) ([]Account, error) {
	return s.accountRepo.FindAll(ctx, userID)
}

func (s *service) FindByID(ctx context.Context, id, userID string) (*Account, error) {
	category, err := s.accountRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *service) Create(ctx context.Context, userID string, req *CreateAccountRequest) (*Account, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	category := &Account{UserID: userID, Name: req.Name}
	err := s.accountRepo.Insert(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *service) Update(ctx context.Context, id string, userID string, req *UpdateAccountRequest) (*Account, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	category, err := s.accountRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	category.Name = req.Name
	err = s.accountRepo.UpdateByID(ctx, category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found or conflict")
		}
		return nil, err
	}
	return category, nil
}

func (s *service) Delete(ctx context.Context, id, userId string) error {
	err := s.accountRepo.DeleteByID(ctx, id, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("category not found")
		}
		return err
	}
	return nil
}
