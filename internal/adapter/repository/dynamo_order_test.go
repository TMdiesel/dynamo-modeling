package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
	"dynamo-modeling/internal/infrastructure"
)

func TestOrderItemConversion(t *testing.T) {
	// Arrange
	orderID, err := value.NewOrderID("test-order-123")
	require.NoError(t, err)

	customerID, err := value.NewCustomerID("test-customer-123")
	require.NoError(t, err)

	productID1, err := value.NewProductID("product-1")
	require.NoError(t, err)

	productID2, err := value.NewProductID("product-2")
	require.NoError(t, err)

	price1, err := value.NewMoney(1299) // $12.99
	require.NoError(t, err)

	price2, err := value.NewMoney(2499) // $24.99
	require.NoError(t, err)

	orderItem1, err := entity.NewOrderItem(productID1, 2, price1)
	require.NoError(t, err)

	orderItem2, err := entity.NewOrderItem(productID2, 1, price2)
	require.NoError(t, err)

	order, err := entity.NewOrder(orderID, customerID, []entity.OrderItem{*orderItem1, *orderItem2})
	require.NoError(t, err)

	// Act: Convert entity to item
	item, err := OrderItemFromEntity(order)
	require.NoError(t, err)

	// Assert: Check item structure
	assert.Equal(t, "ORDER#test-order-123", item.PK)
	assert.Equal(t, "ORDER#test-order-123", item.SK)
	assert.Equal(t, "CUSTOMER#test-customer-123", item.GSI1PK)
	assert.Contains(t, item.GSI1SK, "ORDER#")
	assert.Contains(t, item.GSI1SK, "test-order-123")
	assert.Equal(t, "ORDER", item.Type)
	assert.Equal(t, "test-order-123", item.ID)
	assert.Equal(t, "test-customer-123", item.CustomerID)
	assert.Equal(t, string(order.Status()), item.Status)
	assert.Equal(t, order.Total().Cents(), item.Total)

	// Check Items JSON contains the products
	assert.Contains(t, item.Items, "product-1")
	assert.Contains(t, item.Items, "product-2")

	// Act: Convert item back to entity
	convertedOrder, err := item.ToEntity()
	require.NoError(t, err)

	// Assert: Check entity structure
	assert.Equal(t, orderID, convertedOrder.ID())
	assert.Equal(t, customerID, convertedOrder.CustomerID())
	assert.Equal(t, order.Status(), convertedOrder.Status())
	assert.Equal(t, order.Total(), convertedOrder.Total())
	assert.Len(t, convertedOrder.Items(), 2)

	items := convertedOrder.Items()
	assert.Equal(t, productID1, items[0].ProductID)
	assert.Equal(t, 2, items[0].Quantity)
	assert.Equal(t, price1, items[0].UnitPrice)

	assert.Equal(t, productID2, items[1].ProductID)
	assert.Equal(t, 1, items[1].Quantity)
	assert.Equal(t, price2, items[1].UnitPrice)
}

// TestDynamoOrderRepository runs integration tests against DynamoDB Local
func TestDynamoOrderRepository(t *testing.T) {
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
	repo := NewDynamoOrderRepository(client)

	t.Run("save and find order", func(t *testing.T) {
		orderID, err := value.NewOrderID("test-save-order")
		require.NoError(t, err)

		customerID, err := value.NewCustomerID("test-customer-order")
		require.NoError(t, err)

		productID, err := value.NewProductID("test-product-order")
		require.NoError(t, err)

		price, err := value.NewMoney(1999)
		require.NoError(t, err)

		orderItem, err := entity.NewOrderItem(productID, 3, price)
		require.NoError(t, err)

		order, err := entity.NewOrder(orderID, customerID, []entity.OrderItem{*orderItem})
		require.NoError(t, err)

		// Save order
		err = repo.Save(ctx, order)
		require.NoError(t, err)

		// Find by ID
		found, err := repo.FindByID(ctx, orderID)
		require.NoError(t, err)
		assert.Equal(t, order.ID(), found.ID())
		assert.Equal(t, order.CustomerID(), found.CustomerID())
		assert.Equal(t, order.Status(), found.Status())
		assert.Equal(t, order.Total(), found.Total())
		assert.Len(t, found.Items(), 1)

		foundItems := found.Items()
		assert.Equal(t, productID, foundItems[0].ProductID)
		assert.Equal(t, 3, foundItems[0].Quantity)
		assert.Equal(t, price, foundItems[0].UnitPrice)

		// Clean up
		err = repo.Delete(ctx, orderID)
		assert.NoError(t, err)
	})

	t.Run("find by customer ID", func(t *testing.T) {
		customerID, err := value.NewCustomerID("test-customer-multi")
		require.NoError(t, err)

		// Create multiple orders for the same customer
		orders := make([]*entity.Order, 2)
		for i := 0; i < 2; i++ {
			orderID, err := value.NewOrderID(fmt.Sprintf("test-multi-order-%d", i))
			require.NoError(t, err)

			productID, err := value.NewProductID(fmt.Sprintf("test-product-multi-%d", i))
			require.NoError(t, err)

			price, err := value.NewMoney(int64(1000 + i*500))
			require.NoError(t, err)

			orderItem, err := entity.NewOrderItem(productID, i+1, price)
			require.NoError(t, err)

			order, err := entity.NewOrder(orderID, customerID, []entity.OrderItem{*orderItem})
			require.NoError(t, err)

			orders[i] = order

			// Save each order
			err = repo.Save(ctx, order)
			require.NoError(t, err)
		}

		// Find by customer ID
		found, err := repo.FindByCustomerID(ctx, customerID)
		require.NoError(t, err)
		assert.Len(t, found, 2)

		// Verify all orders are present
		foundIDs := make([]string, len(found))
		for i, o := range found {
			foundIDs[i] = o.ID().String()
		}

		for _, original := range orders {
			assert.Contains(t, foundIDs, original.ID().String())
		}

		// Clean up
		for _, order := range orders {
			repo.Delete(ctx, order.ID())
		}
	})

	t.Run("find by ID not found", func(t *testing.T) {
		orderID, err := value.NewOrderID("non-existent-order")
		require.NoError(t, err)

		// Act
		order, err := repo.FindByID(ctx, orderID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Contains(t, err.Error(), "order not found")
	})

	t.Run("delete order", func(t *testing.T) {
		orderID, err := value.NewOrderID("test-delete-order")
		require.NoError(t, err)

		customerID, err := value.NewCustomerID("test-customer-delete")
		require.NoError(t, err)

		productID, err := value.NewProductID("test-product-delete")
		require.NoError(t, err)

		price, err := value.NewMoney(799)
		require.NoError(t, err)

		orderItem, err := entity.NewOrderItem(productID, 1, price)
		require.NoError(t, err)

		order, err := entity.NewOrder(orderID, customerID, []entity.OrderItem{*orderItem})
		require.NoError(t, err)

		// Save first
		err = repo.Save(ctx, order)
		require.NoError(t, err)

		// Verify it exists
		found, err := repo.FindByID(ctx, orderID)
		require.NoError(t, err)
		assert.NotNil(t, found)

		// Act: Delete
		err = repo.Delete(ctx, orderID)
		require.NoError(t, err)

		// Assert: Verify it's gone
		deleted, err := repo.FindByID(ctx, orderID)
		assert.Error(t, err)
		assert.Nil(t, deleted)
		assert.Contains(t, err.Error(), "order not found")
	})
}
