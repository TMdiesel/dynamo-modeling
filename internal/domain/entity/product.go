package entity

import (
	"fmt"
	"time"

	"dynamo-modeling/internal/domain/value"
)

// Product represents a product entity
type Product struct {
	id          value.ProductID
	name        string
	description string
	price       value.Money
	stock       int
	createdAt   time.Time
	updatedAt   time.Time
}

// NewProduct creates a new Product entity
func NewProduct(id value.ProductID, name, description string, price value.Money, stock int) (*Product, error) {
	if name == "" {
		return nil, fmt.Errorf("product name cannot be empty")
	}
	if stock < 0 {
		return nil, fmt.Errorf("product stock cannot be negative")
	}

	now := time.Now()
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ID returns the product ID
func (p *Product) ID() value.ProductID {
	return p.id
}

// Name returns the product name
func (p *Product) Name() string {
	return p.name
}

// Description returns the product description
func (p *Product) Description() string {
	return p.description
}

// Price returns the product price
func (p *Product) Price() value.Money {
	return p.price
}

// Stock returns the current stock level
func (p *Product) Stock() int {
	return p.stock
}

// CreatedAt returns the creation timestamp
func (p *Product) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the last update timestamp
func (p *Product) UpdatedAt() time.Time {
	return p.updatedAt
}

// UpdatePrice updates the product price
func (p *Product) UpdatePrice(price value.Money) {
	p.price = price
	p.updatedAt = time.Now()
}

// UpdateStock sets the stock level
func (p *Product) UpdateStock(stock int) error {
	if stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}
	p.stock = stock
	p.updatedAt = time.Now()
	return nil
}

// AddStock increases the stock level
func (p *Product) AddStock(amount int) error {
	if amount < 0 {
		return fmt.Errorf("amount to add cannot be negative")
	}
	p.stock += amount
	p.updatedAt = time.Now()
	return nil
}

// ReserveStock decreases the stock level (for orders)
func (p *Product) ReserveStock(amount int) error {
	if amount < 0 {
		return fmt.Errorf("amount to reserve cannot be negative")
	}
	if p.stock < amount {
		return fmt.Errorf("insufficient stock: available %d, requested %d", p.stock, amount)
	}
	p.stock -= amount
	p.updatedAt = time.Now()
	return nil
}

// IsInStock checks if the product has sufficient stock
func (p *Product) IsInStock(quantity int) bool {
	return p.stock >= quantity
}

// IsAvailable checks if the product is available for purchase
func (p *Product) IsAvailable() bool {
	return p.stock > 0
}

// UpdateDetails updates the product name and description
func (p *Product) UpdateDetails(name, description string) error {
	if name == "" {
		return fmt.Errorf("product name cannot be empty")
	}
	p.name = name
	p.description = description
	p.updatedAt = time.Now()
	return nil
}

// Equals compares two Product entities by ID
func (p *Product) Equals(other *Product) bool {
	if other == nil {
		return false
	}
	return p.id == other.id
}
