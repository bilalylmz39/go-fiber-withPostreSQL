package product

import (
	"context"
	"errors"
	"fmt"

	domain "ybilaly/internal/domain/product"
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, input domain.CreateInput) (domain.Product, error) {
	product, err := domain.New(input)
	if err != nil {
		return domain.Product{}, err
	}

	return s.repository.Create(ctx, product)
}

func (s *Service) Update(ctx context.Context, id int64, input domain.UpdateInput) (domain.Product, error) {
	if id < 1 {
		return domain.Product{}, fmt.Errorf("%w: invalid product id", domain.ErrValidation)
	}

	current, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}

	if err := current.ApplyUpdate(input); err != nil {
		return domain.Product{}, err
	}

	return s.repository.Update(ctx, current)
}

func (s *Service) GetByID(ctx context.Context, id int64) (domain.Product, error) {
	if id < 1 {
		return domain.Product{}, fmt.Errorf("%w: invalid product id", domain.ErrValidation)
	}

	return s.repository.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter domain.ListFilter) ([]domain.Product, error) {
	filter = normalizeFilter(filter)
	return s.repository.List(ctx, filter)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return fmt.Errorf("%w: invalid product id", domain.ErrValidation)
	}

	return s.repository.Delete(ctx, id)
}

func (s *Service) ChangeStock(ctx context.Context, id int64, delta int) (domain.Product, error) {
	current, err := s.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}

	if err := current.ChangeStock(delta); err != nil {
		return domain.Product{}, err
	}

	return s.repository.Update(ctx, current)
}

func (s *Service) SetActive(ctx context.Context, id int64, active bool) (domain.Product, error) {
	current, err := s.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}

	if active {
		current.Activate()
	} else {
		current.Deactivate()
	}

	return s.repository.Update(ctx, current)
}

func IsValidationError(err error) bool {
	return errors.Is(err, domain.ErrValidation)
}

func IsNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

func normalizeFilter(filter domain.ListFilter) domain.ListFilter {
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return filter
}
