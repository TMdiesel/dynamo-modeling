package repository

import (
	"context"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
)

// CustomerRepository defines the interface for customer persistence operations
type CustomerRepository interface {
	// Save creates or updates a customer
	Save(ctx context.Context, customer *entity.Customer) error

	// FindByID retrieves a customer by their ID
	FindByID(ctx context.Context, id value.CustomerID) (*entity.Customer, error)

	// FindByEmail retrieves a customer by their email address
	FindByEmail(ctx context.Context, email value.Email) (*entity.Customer, error)

	// Delete removes a customer by their ID
	Delete(ctx context.Context, id value.CustomerID) error

	// Exists checks if a customer exists by their ID
	Exists(ctx context.Context, id value.CustomerID) (bool, error)

	// ListWithLimit retrieves customers with optional limit
	ListWithLimit(ctx context.Context, limit *int) ([]*entity.Customer, error)
}
