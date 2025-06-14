package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerID(t *testing.T) {
	t.Run("valid customer ID", func(t *testing.T) {
		id, err := NewCustomerID("customer-123")
		assert.NoError(t, err)
		assert.Equal(t, "customer-123", id.String())
		assert.False(t, id.IsEmpty())
	})

	t.Run("empty customer ID should return error", func(t *testing.T) {
		_, err := NewCustomerID("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("whitespace only customer ID should return error", func(t *testing.T) {
		_, err := NewCustomerID("   ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("empty customer ID check", func(t *testing.T) {
		var id CustomerID
		assert.True(t, id.IsEmpty())
	})
}

func TestProductID(t *testing.T) {
	t.Run("valid product ID", func(t *testing.T) {
		id, err := NewProductID("product-456")
		assert.NoError(t, err)
		assert.Equal(t, "product-456", id.String())
		assert.False(t, id.IsEmpty())
	})

	t.Run("empty product ID should return error", func(t *testing.T) {
		_, err := NewProductID("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("whitespace only product ID should return error", func(t *testing.T) {
		_, err := NewProductID("   ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("empty product ID check", func(t *testing.T) {
		var id ProductID
		assert.True(t, id.IsEmpty())
	})
}

func TestOrderID(t *testing.T) {
	t.Run("valid order ID", func(t *testing.T) {
		id, err := NewOrderID("order-789")
		assert.NoError(t, err)
		assert.Equal(t, "order-789", id.String())
		assert.False(t, id.IsEmpty())
	})

	t.Run("empty order ID should return error", func(t *testing.T) {
		_, err := NewOrderID("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("whitespace only order ID should return error", func(t *testing.T) {
		_, err := NewOrderID("   ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("empty order ID check", func(t *testing.T) {
		var id OrderID
		assert.True(t, id.IsEmpty())
	})
}
