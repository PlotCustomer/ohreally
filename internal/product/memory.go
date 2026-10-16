package product

import (
	"context"
	"sync"
)

// MemoryRepository keeps the Product DataObject in memory. It backs the
// Webshop when no Postgres database is available, so the delivery can always
// be previewed.
type MemoryRepository struct {
	mu       sync.RWMutex
	products map[string]Product
	order    []string
}

// NewMemoryRepository returns a repository seeded with the given products.
func NewMemoryRepository(seed []Product) *MemoryRepository {
	repo := &MemoryRepository{products: make(map[string]Product, len(seed))}
	for _, p := range seed {
		if _, exists := repo.products[p.ID]; !exists {
			repo.order = append(repo.order, p.ID)
		}
		repo.products[p.ID] = p
	}
	return repo
}

// List returns all products in catalog order.
func (r *MemoryRepository) List(_ context.Context) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]Product, 0, len(r.order))
	for _, id := range r.order {
		products = append(products, r.products[id])
	}
	return products, nil
}

// Get returns the product with the given id.
func (r *MemoryRepository) Get(_ context.Context, id string) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

// Close releases the resources held by the repository.
func (r *MemoryRepository) Close() error { return nil }
