// Package selection models the Selecteer producten business process: the Klant
// selects Products from the Webshop catalog.
package selection

import (
	"sync"

	"dvbd.com/webshop/internal/product"
)

// Line is one selected product with its quantity.
type Line struct {
	Product  product.Product
	Quantity int
}

// Store keeps the selection of every Klant in memory.
type Store struct {
	mu      sync.Mutex
	byKlant map[string]map[string]int
}

// NewStore returns an empty selection store.
func NewStore() *Store {
	return &Store{byKlant: make(map[string]map[string]int)}
}

// Add selects quantity of the given product for the Klant.
func (s *Store) Add(klantID string, productID string, quantity int) {
	if quantity <= 0 {
		quantity = 1
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	lines := s.byKlant[klantID]
	if lines == nil {
		lines = make(map[string]int)
		s.byKlant[klantID] = lines
	}
	lines[productID] += quantity
}

// Remove drops the product from the selection of the Klant.
func (s *Store) Remove(klantID string, productID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if lines := s.byKlant[klantID]; lines != nil {
		delete(lines, productID)
		if len(lines) == 0 {
			delete(s.byKlant, klantID)
		}
	}
}

// Clear empties the selection of the Klant.
func (s *Store) Clear(klantID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.byKlant, klantID)
}

// Count returns the number of selected products of the Klant.
func (s *Store) Count(klantID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	total := 0
	for _, quantity := range s.byKlant[klantID] {
		total += quantity
	}
	return total
}

// Lines returns the selection of the Klant, resolved against the catalog and
// ordered by the catalog order of the products.
func (s *Store) Lines(klantID string, catalog []product.Product) []Line {
	s.mu.Lock()
	quantities := make(map[string]int, len(s.byKlant[klantID]))
	for productID, quantity := range s.byKlant[klantID] {
		quantities[productID] = quantity
	}
	s.mu.Unlock()

	lines := make([]Line, 0, len(quantities))
	for _, p := range catalog {
		if quantity, ok := quantities[p.ID]; ok {
			lines = append(lines, Line{Product: p, Quantity: quantity})
			delete(quantities, p.ID)
		}
	}
	// Include selected products that are no longer in the catalog so the Klant
	// can still see and remove them.
	for productID, quantity := range quantities {
		lines = append(lines, Line{Product: product.Product{ID: productID, Name: productID, Category: product.Category("Onbekend")}, Quantity: quantity})
	}
	return lines
}
