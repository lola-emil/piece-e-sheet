package category

import (
	"context"
	"database/sql"
	"errors"
)

type service struct {
	categoryRepo CategoryRepository
}

type CategoryService interface {
	FindAll(ctx context.Context, userID string) ([]Category, error)
	FindByID(ctx context.Context, id string) (*Category, error)
	Create(ctx context.Context, userID string, req *CreateCategoryRequest) (*Category, error)
	Update(ctx context.Context, id string, req *UpdateCategoryRequest) (*Category, error)
	Delete(ctx context.Context, id string) error
}

func NewCategoryService(categoryRepo CategoryRepository) CategoryService {
	return &service{categoryRepo: categoryRepo}
}

func (s *service) FindAll(ctx context.Context, userID string) ([]Category, error) {
	return s.categoryRepo.FindAll(ctx, userID)
}

func (s *service) FindByID(ctx context.Context, id string) (*Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *service) Create(ctx context.Context, userID string, req *CreateCategoryRequest) (*Category, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	category := &Category{UserID: userID, Name: req.Name}
	err := s.categoryRepo.Insert(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *service) Update(ctx context.Context, id string, req *UpdateCategoryRequest) (*Category, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	category.Name = req.Name
	err = s.categoryRepo.UpdateByID(ctx, category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found or conflict")
		}
		return nil, err
	}
	return category, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.categoryRepo.DeleteByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("category not found")
		}
		return err
	}
	return nil
}
