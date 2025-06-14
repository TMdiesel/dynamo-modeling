package repository

import (
	"context"
	"fmt"
	"sync"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
)

// InMemoryProductRepository is an in-memory implementation of ProductRepository for testing
type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[string]*entity.Product
}

// NewInMemoryProductRepository creates a new in-memory product repository
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*entity.Product),
	}
}

// Save creates or updates a product
func (r *InMemoryProductRepository) Save(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[product.ID().String()] = product
	return nil
}

// FindByID retrieves a product by its ID
func (r *InMemoryProductRepository) FindByID(ctx context.Context, id value.ProductID) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id.String()]
	if !exists {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

// FindAll retrieves all products with optional pagination
func (r *InMemoryProductRepository) FindAll(ctx context.Context, limit int, lastKey *string) ([]*entity.Product, *string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var products []*entity.Product
	var startKey string
	if lastKey != nil {
		startKey = *lastKey
	}

	startFound := startKey == ""
	count := 0

	for id, product := range r.products {
		if !startFound {
			if id == startKey {
				startFound = true
			}
			continue
		}

		if count >= limit {
			// Return the next key for pagination
			nextKey := id
			return products, &nextKey, nil
		}

		products = append(products, product)
		count++
	}

	return products, nil, nil
}

// FindInStock retrieves products that are currently in stock
func (r *InMemoryProductRepository) FindInStock(ctx context.Context, limit int, lastKey *string) ([]*entity.Product, *string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var products []*entity.Product
	var startKey string
	if lastKey != nil {
		startKey = *lastKey
	}

	startFound := startKey == ""
	count := 0

	for id, product := range r.products {
		if !product.IsAvailable() {
			continue
		}

		if !startFound {
			if id == startKey {
				startFound = true
			}
			continue
		}

		if count >= limit {
			// Return the next key for pagination
			nextKey := id
			return products, &nextKey, nil
		}

		products = append(products, product)
		count++
	}

	return products, nil, nil
}

// Delete removes a product by its ID
func (r *InMemoryProductRepository) Delete(ctx context.Context, id value.ProductID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id.String()]; !exists {
		return fmt.Errorf("product not found")
	}

	delete(r.products, id.String())
	return nil
}

// Exists checks if a product exists by its ID
func (r *InMemoryProductRepository) Exists(ctx context.Context, id value.ProductID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.products[id.String()]
	return exists, nil
}
