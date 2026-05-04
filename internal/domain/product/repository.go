package product

import "context"

type Repository interface {
	Create(ctx context.Context, product Product) (Product, error)
	Update(ctx context.Context, product Product) (Product, error)
	FindByID(ctx context.Context, id int64) (Product, error)
	List(ctx context.Context, filter ListFilter) ([]Product, error)
	Delete(ctx context.Context, id int64) error
}
