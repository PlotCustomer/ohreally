// Package product models the Product DataObject of the Webshop, together with
// the BusinessObject specializations it realises: Fruit & Groenten and Pasta.
package product

import (
	"context"
	"errors"
)

// Category is a specialization of the Product BusinessObject. The Webshop
// offers products in the categories Fruit & Groenten and Pasta.
type Category string

const (
	// CategoryFruitGroenten is the Fruit & Groenten specialization.
	CategoryFruitGroenten Category = "Fruit & Groenten"
	// CategoryPasta is the Pasta specialization.
	CategoryPasta Category = "Pasta"
)

// Categories lists the Product specializations in catalog order.
var Categories = []Category{CategoryFruitGroenten, CategoryPasta}

// Product is the Product DataObject: an article offered by the Webshop.
type Product struct {
	ID          string
	Name        string
	Category    Category
	Description string
	PriceCents  int64
	Currency    string
	Unit        string
	Stock       int
}

// ErrNotFound is returned when a product is not part of the catalog.
var ErrNotFound = errors.New("product not found")

// Repository provides access to the Product DataObject.
type Repository interface {
	// List returns all products of the catalog.
	List(ctx context.Context) ([]Product, error)
	// Get returns the product with the given id.
	Get(ctx context.Context, id string) (Product, error)
	// Close releases the resources held by the repository.
	Close() error
}
