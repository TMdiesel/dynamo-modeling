//go:build integration
// +build integration

package repository

import (
"context"
"testing"
"time"

"github.com/google/uuid"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"

"dynamo-modeling/internal/domain/entity"
"dynamo-modeling/internal/domain/value"
"dynamo-modeling/internal/infrastructure"
)

// TestRepositoriesIntegration tests all repositories working together
func TestRepositoriesIntegration(t *testing.T) {
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

	// Create repositories
	customerRepo := NewDynamoCustomerRepository(client)
	productRepo := NewDynamoProductRepository(client)
	orderRepo := NewDynamoOrderRepository(client)

	t.Run("end-to-end e-commerce scenario", func(t *testing.T) {
// Create test customer
customerID, _ := value.NewCustomerID(uuid.New().String())
		email, _ := value.NewEmail("integration@example.com")
		customer := entity.NewCustomer(customerID, email, "Integration Test Customer")

		// Save customer
		err := customerRepo.Save(ctx, customer)
		require.NoError(t, err)

		// Create test product
		productID, _ := value.NewProductID(uuid.New().String())
		money, _ := value.NewMoney(1500)
		product, err := entity.NewProduct(productID, "Integration Test Product", "Test Description", money, 100)
		require.NoError(t, err)

		// Save product
		err = productRepo.Save(ctx, product)
		require.NoError(t, err)

		// Create test order
		orderID, _ := value.NewOrderID(uuid.New().String())
		orderItem, err := entity.NewOrderItem(productID, 2, money)
		require.NoError(t, err)
		order, err := entity.NewOrder(orderID, customerID, []entity.OrderItem{*orderItem})
		require.NoError(t, err)

		// Save order
		err = orderRepo.Save(ctx, order)
		require.NoError(t, err)

		// Verify all data can be retrieved
		retrievedCustomer, err := customerRepo.FindByID(ctx, customerID)
		require.NoError(t, err)
		assert.Equal(t, customer.ID(), retrievedCustomer.ID())
		assert.Equal(t, customer.Email(), retrievedCustomer.Email())
		assert.Equal(t, customer.Name(), retrievedCustomer.Name())

		retrievedProduct, err := productRepo.FindByID(ctx, productID)
		require.NoError(t, err)
		assert.Equal(t, product.ID(), retrievedProduct.ID())
		assert.Equal(t, product.Name(), retrievedProduct.Name())
		assert.Equal(t, product.Price(), retrievedProduct.Price())

		retrievedOrder, err := orderRepo.FindByID(ctx, orderID)
		require.NoError(t, err)
		assert.Equal(t, order.ID(), retrievedOrder.ID())
		assert.Equal(t, order.CustomerID(), retrievedOrder.CustomerID())
		assert.Len(t, retrievedOrder.Items(), 1)
		assert.Equal(t, orderItem.ProductID, retrievedOrder.Items()[0].ProductID)
		assert.Equal(t, orderItem.Quantity, retrievedOrder.Items()[0].Quantity)
	})

	t.Run("bulk operations performance test", func(t *testing.T) {
const numRecords = 20 // Reduced for faster testing
start := time.Now()

		// Create multiple customers
		var customers []*entity.Customer
		for i := 0; i < numRecords; i++ {
			customerID, _ := value.NewCustomerID(uuid.New().String())
			email, _ := value.NewEmail("bulk-" + uuid.New().String() + "@example.com")
			customer := entity.NewCustomer(customerID, email, "Bulk Customer")
			customers = append(customers, customer)

			err := customerRepo.Save(ctx, customer)
			require.NoError(t, err)
		}

		// Create multiple products
		var products []*entity.Product
		for i := 0; i < numRecords; i++ {
			productID, _ := value.NewProductID(uuid.New().String())
			money, _ := value.NewMoney(int64(1000 + i))
			product, err := entity.NewProduct(productID, "Bulk Product", "Description", money, 50)
			require.NoError(t, err)
			products = append(products, product)

			err = productRepo.Save(ctx, product)
			require.NoError(t, err)
		}

		// Create multiple orders
		for i := 0; i < numRecords; i++ {
			orderID, _ := value.NewOrderID(uuid.New().String())
			orderItem, err := entity.NewOrderItem(products[i].ID(), 1, products[i].Price())
			require.NoError(t, err)
			order, err := entity.NewOrder(orderID, customers[i].ID(), []entity.OrderItem{*orderItem})
			require.NoError(t, err)

			err = orderRepo.Save(ctx, order)
			require.NoError(t, err)
		}

		elapsed := time.Since(start)
		t.Logf("Bulk operations (%d records each for customers, products, orders) completed in %v", numRecords, elapsed)

		// Verify data integrity by checking some records
		retrievedCustomer, err := customerRepo.FindByID(ctx, customers[0].ID())
		require.NoError(t, err)
		assert.Equal(t, customers[0].ID(), retrievedCustomer.ID())

		retrievedProduct, err := productRepo.FindByID(ctx, products[0].ID())
		require.NoError(t, err)
		assert.Equal(t, products[0].ID(), retrievedProduct.ID())
	})
}
