package repository

import (
	"context"
	"fmt"
	"sync"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
)

// InMemoryCustomerRepository is an in-memory implementation of CustomerRepository for testing
type InMemoryCustomerRepository struct {
	mu        sync.RWMutex
	customers map[string]*entity.Customer
	emails    map[string]value.CustomerID // email -> customerID mapping
}

// NewInMemoryCustomerRepository creates a new in-memory customer repository
func NewInMemoryCustomerRepository() *InMemoryCustomerRepository {
	return &InMemoryCustomerRepository{
		customers: make(map[string]*entity.Customer),
		emails:    make(map[string]value.CustomerID),
	}
}

// Save creates or updates a customer
func (r *InMemoryCustomerRepository) Save(ctx context.Context, customer *entity.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email is already taken by another customer
	if existingID, exists := r.emails[customer.Email().String()]; exists && existingID != customer.ID() {
		return fmt.Errorf("email already taken by another customer")
	}

	r.customers[customer.ID().String()] = customer
	r.emails[customer.Email().String()] = customer.ID()
	return nil
}

// FindByID retrieves a customer by their ID
func (r *InMemoryCustomerRepository) FindByID(ctx context.Context, id value.CustomerID) (*entity.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customer, exists := r.customers[id.String()]
	if !exists {
		return nil, fmt.Errorf("customer not found")
	}
	return customer, nil
}

// FindByEmail retrieves a customer by their email address
func (r *InMemoryCustomerRepository) FindByEmail(ctx context.Context, email value.Email) (*entity.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	customerID, exists := r.emails[email.String()]
	if !exists {
		return nil, fmt.Errorf("customer not found")
	}

	customer, exists := r.customers[customerID.String()]
	if !exists {
		return nil, fmt.Errorf("customer not found")
	}
	return customer, nil
}

// Delete removes a customer by their ID
func (r *InMemoryCustomerRepository) Delete(ctx context.Context, id value.CustomerID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	customer, exists := r.customers[id.String()]
	if !exists {
		return fmt.Errorf("customer not found")
	}

	delete(r.customers, id.String())
	delete(r.emails, customer.Email().String())
	return nil
}

// Exists checks if a customer exists by their ID
func (r *InMemoryCustomerRepository) Exists(ctx context.Context, id value.CustomerID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.customers[id.String()]
	return exists, nil
}
