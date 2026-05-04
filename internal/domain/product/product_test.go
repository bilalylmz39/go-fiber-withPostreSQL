package product

import (
	"errors"
	"testing"
)

func TestNewNormalizesSKUAndTitle(t *testing.T) {
	product, err := New(CreateInput{
		SKU:   " bk-001 ",
		Title: " Go Book ",
		Price: 10,
		Stock: 5,
	})

	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if product.SKU != "BK-001" {
		t.Fatalf("SKU = %q, want BK-001", product.SKU)
	}
	if product.Title != "Go Book" {
		t.Fatalf("Title = %q, want Go Book", product.Title)
	}
	if !product.Active {
		t.Fatal("new product should be active")
	}
}

func TestChangeStockRejectsNegativeStock(t *testing.T) {
	product, err := New(CreateInput{
		SKU:   "BK-001",
		Title: "Go Book",
		Price: 10,
		Stock: 2,
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	err = product.ChangeStock(-3)
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("err = %v, want validation error", err)
	}
	if product.Stock != 2 {
		t.Fatalf("stock = %d, want 2", product.Stock)
	}
}

func TestValidateRejectsMissingSKU(t *testing.T) {
	_, err := New(CreateInput{
		Title: "Go Book",
		Price: 10,
		Stock: 2,
	})

	if !errors.Is(err, ErrValidation) {
		t.Fatalf("err = %v, want validation error", err)
	}
}
