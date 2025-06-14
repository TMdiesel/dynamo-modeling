package value

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

// CustomerID represents a unique customer identifier
type CustomerID string

// NewCustomerID creates a new CustomerID with validation
func NewCustomerID(id string) (CustomerID, error) {
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("customer ID cannot be empty")
	}
	return CustomerID(id), nil
}

// String returns the string representation of CustomerID
func (c CustomerID) String() string {
	return string(c)
}

// IsEmpty checks if the CustomerID is empty
func (c CustomerID) IsEmpty() bool {
	return string(c) == ""
}

// ProductID represents a unique product identifier
type ProductID string

// NewProductID creates a new ProductID with validation
func NewProductID(id string) (ProductID, error) {
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("product ID cannot be empty")
	}
	return ProductID(id), nil
}

// String returns the string representation of ProductID
func (p ProductID) String() string {
	return string(p)
}

// IsEmpty checks if the ProductID is empty
func (p ProductID) IsEmpty() bool {
	return string(p) == ""
}

// OrderID represents a unique order identifier
type OrderID string

// NewOrderID creates a new OrderID with validation
func NewOrderID(id string) (OrderID, error) {
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("order ID cannot be empty")
	}
	return OrderID(id), nil
}

// String returns the string representation of OrderID
func (o OrderID) String() string {
	return string(o)
}

// IsEmpty checks if the OrderID is empty
func (o OrderID) IsEmpty() bool {
	return string(o) == ""
}

// GenerateCustomerID generates a new unique CustomerID
func GenerateCustomerID() CustomerID {
	id := generateUUID()
	return CustomerID(id)
}

// GenerateProductID generates a new unique ProductID
func GenerateProductID() ProductID {
	id := generateUUID()
	return ProductID(id)
}

// GenerateOrderID generates a new unique OrderID
func GenerateOrderID() OrderID {
	id := generateUUID()
	return OrderID(id)
}

// generateUUID generates a simple UUID v4
func generateUUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}

	// Set version (4) and variant bits
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
