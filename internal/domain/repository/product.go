package repository

import (
	"context"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
)

// ProductRepository defines the interface for product persistence operations
type ProductRepository interface {
	// Save creates or updates a product
	Save(ctx context.Context, product *entity.Product) error

	// FindByID retrieves a product by its ID
	FindByID(ctx context.Context, id value.ProductID) (*entity.Product, error)

	// FindAll retrieves all products with optional pagination
	FindAll(ctx context.Context, limit int, lastKey *string) ([]*entity.Product, *string, error)

	// FindInStock retrieves products that are currently in stock
	FindInStock(ctx context.Context, limit int, lastKey *string) ([]*entity.Product, *string, error)

	// Delete removes a product by its ID
	Delete(ctx context.Context, id value.ProductID) error

	// Exists checks if a product exists by its ID
	Exists(ctx context.Context, id value.ProductID) (bool, error)
}
