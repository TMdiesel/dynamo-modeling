package entity

import (
	"time"

	"dynamo-modeling/internal/domain/value"
)

// Customer represents a customer entity
type Customer struct {
	id        value.CustomerID
	email     value.Email
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewCustomer creates a new Customer entity
func NewCustomer(id value.CustomerID, email value.Email, name string) *Customer {
	now := time.Now()
	return &Customer{
		id:        id,
		email:     email,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the customer ID
func (c *Customer) ID() value.CustomerID {
	return c.id
}

// Email returns the customer email
func (c *Customer) Email() value.Email {
	return c.email
}

// Name returns the customer name
func (c *Customer) Name() string {
	return c.name
}

// CreatedAt returns the creation timestamp
func (c *Customer) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the last update timestamp
func (c *Customer) UpdatedAt() time.Time {
	return c.updatedAt
}

// UpdateEmail updates the customer's email address
func (c *Customer) UpdateEmail(email value.Email) {
	c.email = email
	c.updatedAt = time.Now()
}

// UpdateName updates the customer's name
func (c *Customer) UpdateName(name string) {
	c.name = name
	c.updatedAt = time.Now()
}

// Equals compares two Customer entities by ID
func (c *Customer) Equals(other *Customer) bool {
	if other == nil {
		return false
	}
	return c.id == other.id
}
