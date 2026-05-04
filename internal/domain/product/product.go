package product

import (
	"fmt"
	"strings"
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	SKU         string    `json:"sku"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInput struct {
	SKU         string  `json:"sku"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type UpdateInput struct {
	SKU         string  `json:"sku"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Active      bool    `json:"active"`
}

type ListFilter struct {
	Query    string
	Active   *bool
	MinPrice *float64
	MaxPrice *float64
	Limit    int
	Offset   int
}

func New(input CreateInput) (Product, error) {
	now := time.Now().UTC()
	product := Product{
		SKU:         normalizeSKU(input.SKU),
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Price:       input.Price,
		Stock:       input.Stock,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := product.Validate(); err != nil {
		return Product{}, err
	}

	return product, nil
}

func (p *Product) ApplyUpdate(input UpdateInput) error {
	p.SKU = normalizeSKU(input.SKU)
	p.Title = strings.TrimSpace(input.Title)
	p.Description = strings.TrimSpace(input.Description)
	p.Price = input.Price
	p.Stock = input.Stock
	p.Active = input.Active
	p.UpdatedAt = time.Now().UTC()

	return p.Validate()
}

func (p *Product) ChangeStock(delta int) error {
	nextStock := p.Stock + delta
	if nextStock < 0 {
		return fmt.Errorf("%w: stock cannot be negative", ErrValidation)
	}

	p.Stock = nextStock
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Product) Activate() {
	p.Active = true
	p.UpdatedAt = time.Now().UTC()
}

func (p *Product) Deactivate() {
	p.Active = false
	p.UpdatedAt = time.Now().UTC()
}

func (p Product) Validate() error {
	if p.SKU == "" {
		return fmt.Errorf("%w: sku is required", ErrValidation)
	}
	if p.Title == "" {
		return fmt.Errorf("%w: title is required", ErrValidation)
	}
	if p.Price < 0 {
		return fmt.Errorf("%w: price cannot be negative", ErrValidation)
	}
	if p.Stock < 0 {
		return fmt.Errorf("%w: stock cannot be negative", ErrValidation)
	}

	return nil
}

func normalizeSKU(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}
