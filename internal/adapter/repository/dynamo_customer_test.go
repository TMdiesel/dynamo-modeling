package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
	"dynamo-modeling/internal/infrastructure"
)

// TestDynamoCustomerRepository runs integration tests against DynamoDB Local
func TestDynamoCustomerRepository(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Initialize DynamoDB client for testing
	client, err := infrastructure.NewDynamoDBClient(ctx, infrastructure.DynamoDBConfig{
		Region:    "ap-northeast-1",
		Endpoint:  "http://localhost:8000",
		TableName: "OnlineShop",
	})
	require.NoError(t, err)

	// Health check
	err = client.HealthCheck(ctx)
	require.NoError(t, err)

	// Create repository
	repo := NewDynamoCustomerRepository(client)

	t.Run("save and find customer", func(t *testing.T) {
		customerID, _ := value.NewCustomerID("test-customer-123")
		email, _ := value.NewEmail("test@example.com")
		customer := entity.NewCustomer(customerID, email, "Test Customer")

		// Save customer
		err := repo.Save(ctx, customer)
		assert.NoError(t, err)

		// Find by ID
		found, err := repo.FindByID(ctx, customerID)
		assert.NoError(t, err)
		assert.Equal(t, customer.ID(), found.ID())
		assert.Equal(t, customer.Email(), found.Email())
		assert.Equal(t, customer.Name(), found.Name())

		// Find by email
		found, err = repo.FindByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, customer.ID(), found.ID())

		// Clean up
		err = repo.Delete(ctx, customerID)
		assert.NoError(t, err)
	})

	t.Run("customer not found", func(t *testing.T) {
		nonExistentID, _ := value.NewCustomerID("non-existent-123")
		_, err := repo.FindByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		nonExistentEmail, _ := value.NewEmail("non@existent.com")
		_, err = repo.FindByEmail(ctx, nonExistentEmail)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("duplicate email should fail", func(t *testing.T) {
		email, _ := value.NewEmail("duplicate-test@example.com")
		customerID1, _ := value.NewCustomerID("customer-1-123")
		customerID2, _ := value.NewCustomerID("customer-2-123")

		customer1 := entity.NewCustomer(customerID1, email, "Customer 1")
		customer2 := entity.NewCustomer(customerID2, email, "Customer 2")

		// Save first customer
		err := repo.Save(ctx, customer1)
		assert.NoError(t, err)

		// Try to save second customer with same email
		err = repo.Save(ctx, customer2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already taken")

		// Clean up
		repo.Delete(ctx, customerID1)
	})

	t.Run("exists check", func(t *testing.T) {
		customerID, _ := value.NewCustomerID("exists-test-123")
		email, _ := value.NewEmail("exists-test@example.com")
		customer := entity.NewCustomer(customerID, email, "Exists Test")

		// Check non-existent customer
		exists, err := repo.Exists(ctx, customerID)
		assert.NoError(t, err)
		assert.False(t, exists)

		// Save customer
		err = repo.Save(ctx, customer)
		assert.NoError(t, err)

		// Check existing customer
		exists, err = repo.Exists(ctx, customerID)
		assert.NoError(t, err)
		assert.True(t, exists)

		// Clean up
		err = repo.Delete(ctx, customerID)
		assert.NoError(t, err)
	})
}
