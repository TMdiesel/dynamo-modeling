package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dynamo-modeling/internal/domain/value"
)

func TestProduct(t *testing.T) {
	// Setup test data
	productID, _ := value.NewProductID("product-123")
	name := "Test Product"
	description := "A test product"
	price, _ := value.NewMoney(1000) // $10.00
	stock := 50

	t.Run("create new product", func(t *testing.T) {
		product, err := NewProduct(productID, name, description, price, stock)

		assert.NoError(t, err)
		assert.Equal(t, productID, product.ID())
		assert.Equal(t, name, product.Name())
		assert.Equal(t, description, product.Description())
		assert.Equal(t, price, product.Price())
		assert.Equal(t, stock, product.Stock())
		assert.False(t, product.CreatedAt().IsZero())
		assert.False(t, product.UpdatedAt().IsZero())
	})

	t.Run("create product with empty name should fail", func(t *testing.T) {
		_, err := NewProduct(productID, "", description, price, stock)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("create product with negative stock should fail", func(t *testing.T) {
		_, err := NewProduct(productID, name, description, price, -1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stock cannot be negative")
	})

	t.Run("update product price", func(t *testing.T) {
		product, _ := NewProduct(productID, name, description, price, stock)
		newPrice, _ := value.NewMoney(1500) // $15.00

		product.UpdatePrice(newPrice)

		assert.Equal(t, newPrice, product.Price())
	})

	t.Run("update stock", func(t *testing.T) {
		product, _ := NewProduct(productID, name, description, price, stock)

		err := product.UpdateStock(75)
		assert.NoError(t, err)
		assert.Equal(t, 75, product.Stock())

		err = product.UpdateStock(-1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "stock cannot be negative")
	})

	t.Run("add stock", func(t *testing.T) {
		product, _ := NewProduct(productID, name, description, price, stock)

		err := product.AddStock(25)
		assert.NoError(t, err)
		assert.Equal(t, 75, product.Stock())

		err = product.AddStock(-5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount to add cannot be negative")
	})

	t.Run("reserve stock", func(t *testing.T) {
		product, _ := NewProduct(productID, name, description, price, stock)

		err := product.ReserveStock(20)
		assert.NoError(t, err)
		assert.Equal(t, 30, product.Stock())

		err = product.ReserveStock(40) // More than available
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient stock")

		err = product.ReserveStock(-5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "amount to reserve cannot be negative")
	})

	t.Run("stock availability checks", func(t *testing.T) {
		product, _ := NewProduct(productID, name, description, price, 10)

		assert.True(t, product.IsInStock(5))
		assert.True(t, product.IsInStock(10))
		assert.False(t, product.IsInStock(15))

		assert.True(t, product.IsAvailable())

		product.UpdateStock(0)
		assert.False(t, product.IsAvailable())
	})

	t.Run("update product details", func(t *testing.T) {
		product, _ := NewProduct(productID, name, description, price, stock)

		newName := "Updated Product"
		newDescription := "Updated description"

		err := product.UpdateDetails(newName, newDescription)
		assert.NoError(t, err)
		assert.Equal(t, newName, product.Name())
		assert.Equal(t, newDescription, product.Description())

		err = product.UpdateDetails("", newDescription)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("product equals", func(t *testing.T) {
		product1, _ := NewProduct(productID, name, description, price, stock)
		product2, _ := NewProduct(productID, "Different Name", description, price, stock)

		otherID, _ := value.NewProductID("product-456")
		product3, _ := NewProduct(otherID, name, description, price, stock)

		assert.True(t, product1.Equals(product2))  // Same ID
		assert.False(t, product1.Equals(product3)) // Different ID
		assert.False(t, product1.Equals(nil))      // Nil comparison
	})
}
