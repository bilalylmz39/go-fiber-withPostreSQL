package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	domain "ybilaly/internal/domain/product"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO products (sku, title, description, price, stock, active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, created_at, updated_at`,
		product.SKU,
		product.Title,
		product.Description,
		product.Price,
		product.Stock,
		product.Active,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	return product, err
}

func (r *ProductRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	err := r.db.QueryRowContext(
		ctx,
		`UPDATE products
		 SET sku = $2, title = $3, description = $4, price = $5, stock = $6, active = $7, updated_at = $8
		 WHERE id = $1
		 RETURNING created_at, updated_at`,
		product.ID,
		product.SKU,
		product.Title,
		product.Description,
		product.Price,
		product.Stock,
		product.Active,
		product.UpdatedAt,
	).Scan(&product.CreatedAt, &product.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Product{}, domain.ErrNotFound
	}
	return product, err
}

func (r *ProductRepository) FindByID(ctx context.Context, id int64) (domain.Product, error) {
	product, err := r.scanOne(ctx, `SELECT id, sku, title, description, price, stock, active, created_at, updated_at FROM products WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Product{}, domain.ErrNotFound
	}
	return product, err
}

func (r *ProductRepository) List(ctx context.Context, filter domain.ListFilter) ([]domain.Product, error) {
	query := strings.Builder{}
	query.WriteString(`SELECT id, sku, title, description, price, stock, active, created_at, updated_at FROM products`)

	args := make([]any, 0)
	conditions := make([]string, 0)
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	if filter.Query != "" {
		value := "%" + strings.ToLower(filter.Query) + "%"
		placeholder := addArg(value)
		conditions = append(conditions, fmt.Sprintf("(LOWER(title) LIKE %s OR LOWER(sku) LIKE %s)", placeholder, placeholder))
	}
	if filter.Active != nil {
		conditions = append(conditions, "active = "+addArg(*filter.Active))
	}
	if filter.MinPrice != nil {
		conditions = append(conditions, "price >= "+addArg(*filter.MinPrice))
	}
	if filter.MaxPrice != nil {
		conditions = append(conditions, "price <= "+addArg(*filter.MaxPrice))
	}
	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" ORDER BY id DESC LIMIT ")
	query.WriteString(addArg(filter.Limit))
	query.WriteString(" OFFSET ")
	query.WriteString(addArg(filter.Offset))

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *ProductRepository) scanOne(ctx context.Context, query string, args ...any) (domain.Product, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	var product domain.Product
	err := row.Scan(
		&product.ID,
		&product.SKU,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.Active,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	return product, err
}

type productScanner interface {
	Scan(dest ...any) error
}

func scanProduct(scanner productScanner) (domain.Product, error) {
	var product domain.Product
	err := scanner.Scan(
		&product.ID,
		&product.SKU,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.Active,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	return product, err
}
