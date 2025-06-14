package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
)

func TestInMemoryCustomerRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("save and find customer", func(t *testing.T) {
		repo := NewInMemoryCustomerRepository()

		customerID, _ := value.NewCustomerID("customer-123")
		email, _ := value.NewEmail("test@example.com")
		customer := entity.NewCustomer(customerID, email, "John Doe")

		// Save customer
		err := repo.Save(ctx, customer)
		assert.NoError(t, err)

		// Find by ID
		found, err := repo.FindByID(ctx, customerID)
		assert.NoError(t, err)
		assert.True(t, customer.Equals(found))

		// Find by email
		found, err = repo.FindByEmail(ctx, email)
		assert.NoError(t, err)
		assert.True(t, customer.Equals(found))
	})

	t.Run("customer not found", func(t *testing.T) {
		repo := NewInMemoryCustomerRepository()

		nonExistentID, _ := value.NewCustomerID("non-existent")
		_, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		nonExistentEmail, _ := value.NewEmail("non@existent.com")
		_, err = repo.FindByEmail(ctx, nonExistentEmail)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("duplicate email should fail", func(t *testing.T) {
		repo := NewInMemoryCustomerRepository()

		email, _ := value.NewEmail("duplicate@example.com")
		customerID1, _ := value.NewCustomerID("customer-1")
		customerID2, _ := value.NewCustomerID("customer-2")

		customer1 := entity.NewCustomer(customerID1, email, "Customer 1")
		customer2 := entity.NewCustomer(customerID2, email, "Customer 2")

		// Save first customer
		err := repo.Save(ctx, customer1)
		assert.NoError(t, err)

		// Try to save second customer with same email
		err = repo.Save(ctx, customer2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already taken")
	})

	t.Run("exists check", func(t *testing.T) {
		repo := NewInMemoryCustomerRepository()

		customerID, _ := value.NewCustomerID("customer-123")
		email, _ := value.NewEmail("test@example.com")
		customer := entity.NewCustomer(customerID, email, "John Doe")

		// Check non-existent customer
		exists, err := repo.Exists(ctx, customerID)
		assert.NoError(t, err)
		assert.False(t, exists)

		// Save customer
		repo.Save(ctx, customer)

		// Check existing customer
		exists, err = repo.Exists(ctx, customerID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("delete customer", func(t *testing.T) {
		repo := NewInMemoryCustomerRepository()

		customerID, _ := value.NewCustomerID("customer-123")
		email, _ := value.NewEmail("test@example.com")
		customer := entity.NewCustomer(customerID, email, "John Doe")

		// Save customer
		repo.Save(ctx, customer)

		// Verify customer exists
		exists, _ := repo.Exists(ctx, customerID)
		assert.True(t, exists)

		// Delete customer
		err := repo.Delete(ctx, customerID)
		assert.NoError(t, err)

		// Verify customer no longer exists
		exists, _ = repo.Exists(ctx, customerID)
		assert.False(t, exists)

		// Try to delete non-existent customer
		err = repo.Delete(ctx, customerID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
