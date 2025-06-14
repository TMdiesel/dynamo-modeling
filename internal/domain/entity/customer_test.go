package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"dynamo-modeling/internal/domain/value"
)

func TestCustomer(t *testing.T) {
	// Setup test data
	customerID, _ := value.NewCustomerID("customer-123")
	email, _ := value.NewEmail("test@example.com")
	name := "John Doe"

	t.Run("create new customer", func(t *testing.T) {
		customer := NewCustomer(customerID, email, name)

		assert.Equal(t, customerID, customer.ID())
		assert.Equal(t, email, customer.Email())
		assert.Equal(t, name, customer.Name())
		assert.False(t, customer.CreatedAt().IsZero())
		assert.False(t, customer.UpdatedAt().IsZero())
		assert.Equal(t, customer.CreatedAt(), customer.UpdatedAt())
	})

	t.Run("update customer email", func(t *testing.T) {
		customer := NewCustomer(customerID, email, name)
		originalUpdatedAt := customer.UpdatedAt()

		// Wait a bit to ensure different timestamp
		time.Sleep(1 * time.Millisecond)

		newEmail, _ := value.NewEmail("newemail@example.com")
		customer.UpdateEmail(newEmail)

		assert.Equal(t, newEmail, customer.Email())
		assert.True(t, customer.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("update customer name", func(t *testing.T) {
		customer := NewCustomer(customerID, email, name)
		originalUpdatedAt := customer.UpdatedAt()

		// Wait a bit to ensure different timestamp
		time.Sleep(1 * time.Millisecond)

		newName := "Jane Doe"
		customer.UpdateName(newName)

		assert.Equal(t, newName, customer.Name())
		assert.True(t, customer.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("customer equals", func(t *testing.T) {
		customer1 := NewCustomer(customerID, email, name)
		customer2 := NewCustomer(customerID, email, "Different Name")

		otherID, _ := value.NewCustomerID("customer-456")
		customer3 := NewCustomer(otherID, email, name)

		assert.True(t, customer1.Equals(customer2))  // Same ID
		assert.False(t, customer1.Equals(customer3)) // Different ID
		assert.False(t, customer1.Equals(nil))       // Nil comparison
	})
}
