package category

import "context"

type service struct {
	categoryRepo CategoryRepository
}

type CategoryService interface {
	FindAll(ctx context.Context) ([]Category, error)
	FindByID(ctx context.Context, id string) (*Category, error)
}

func NewCategoryService(categoryRepo CategoryRepository) CategoryService {
	return &service{
		categoryRepo: categoryRepo,
	}
}

func (s *service) FindAll(ctx context.Context) ([]Category, error) {
	return s.categoryRepo.FindAll(ctx)
}

func (s *service) FindByID(ctx context.Context, id string) (*Category, error) {
	return s.categoryRepo.FindByID(ctx, id)
}
