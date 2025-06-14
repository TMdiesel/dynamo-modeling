package repository

import (
	"context"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
)

// OrderRepository defines the interface for order persistence operations
type OrderRepository interface {
	// Save creates or updates an order
	Save(ctx context.Context, order *entity.Order) error

	// FindByID retrieves an order by its ID
	FindByID(ctx context.Context, id value.OrderID) (*entity.Order, error)

	// FindByCustomerID retrieves orders for a specific customer
	FindByCustomerID(ctx context.Context, customerID value.CustomerID, limit int, lastKey *string) ([]*entity.Order, *string, error)

	// FindByStatus retrieves orders with a specific status
	FindByStatus(ctx context.Context, status entity.OrderStatus, limit int, lastKey *string) ([]*entity.Order, *string, error)

	// FindByCustomerAndStatus retrieves orders for a customer with a specific status
	FindByCustomerAndStatus(ctx context.Context, customerID value.CustomerID, status entity.OrderStatus, limit int, lastKey *string) ([]*entity.Order, *string, error)

	// Delete removes an order by its ID
	Delete(ctx context.Context, id value.OrderID) error

	// Exists checks if an order exists by its ID
	Exists(ctx context.Context, id value.OrderID) (bool, error)
}
