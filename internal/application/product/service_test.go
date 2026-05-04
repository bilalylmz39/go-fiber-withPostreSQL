package product

import (
	"context"
	"errors"
	"testing"

	domain "ybilaly/internal/domain/product"
)

func TestServiceCreateNormalizesAndPersistsProduct(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository)

	created, err := service.Create(context.Background(), domain.CreateInput{
		SKU:         " bk-001 ",
		Title:       " Go Book ",
		Description: "Practical Go",
		Price:       42.12,
		Stock:       7,
	})

	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("Create did not assign an id")
	}
	if created.SKU != "BK-001" {
		t.Fatalf("SKU = %q, want %q", created.SKU, "BK-001")
	}
	if !created.Active {
		t.Fatal("new product should be active")
	}
}

func TestServiceCreateRejectsInvalidProduct(t *testing.T) {
	service := NewService(newFakeRepository())

	_, err := service.Create(context.Background(), domain.CreateInput{
		SKU:   "BAD-001",
		Title: "",
		Price: 10,
		Stock: 1,
	})

	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("err = %v, want validation error", err)
	}
}

func TestServiceUpdateReturnsNotFound(t *testing.T) {
	service := NewService(newFakeRepository())

	_, err := service.Update(context.Background(), 999, domain.UpdateInput{
		SKU:    "BK-001",
		Title:  "Go Book",
		Price:  10,
		Stock:  1,
		Active: true,
	})

	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v, want not found", err)
	}
}

func TestServiceChangeStockRejectsNegativeResult(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository)

	created, err := service.Create(context.Background(), domain.CreateInput{
		SKU:   "BK-001",
		Title: "Go Book",
		Price: 10,
		Stock: 3,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	_, err = service.ChangeStock(context.Background(), created.ID, -4)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("err = %v, want validation error", err)
	}

	current, err := service.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if current.Stock != 3 {
		t.Fatalf("stock = %d, want 3", current.Stock)
	}
}

func TestServiceListNormalizesPagination(t *testing.T) {
	repository := newFakeRepository()
	service := NewService(repository)

	_, err := service.List(context.Background(), domain.ListFilter{Limit: 500, Offset: -10})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if repository.lastFilter.Limit != 20 {
		t.Fatalf("limit = %d, want 20", repository.lastFilter.Limit)
	}
	if repository.lastFilter.Offset != 0 {
		t.Fatalf("offset = %d, want 0", repository.lastFilter.Offset)
	}
}

type fakeRepository struct {
	nextID     int64
	products   map[int64]domain.Product
	lastFilter domain.ListFilter
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		nextID:   1,
		products: make(map[int64]domain.Product),
	}
}

func (r *fakeRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	product.ID = r.nextID
	r.nextID++
	r.products[product.ID] = product
	return product, nil
}

func (r *fakeRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	if _, ok := r.products[product.ID]; !ok {
		return domain.Product{}, domain.ErrNotFound
	}
	r.products[product.ID] = product
	return product, nil
}

func (r *fakeRepository) FindByID(ctx context.Context, id int64) (domain.Product, error) {
	product, ok := r.products[id]
	if !ok {
		return domain.Product{}, domain.ErrNotFound
	}
	return product, nil
}

func (r *fakeRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.Product, error) {
	r.lastFilter = filter
	products := make([]domain.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}
	return products, nil
}

func (r *fakeRepository) Delete(ctx context.Context, id int64) error {
	if _, ok := r.products[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.products, id)
	return nil
}
