package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
	"dynamo-modeling/internal/infrastructure"
)

func TestProductItemConversion(t *testing.T) {
	// Arrange
	productID, err := value.NewProductID("test-product-123")
	require.NoError(t, err)

	price, err := value.NewMoney(1299) // $12.99
	require.NoError(t, err)

	product, err := entity.NewProduct(productID, "Test Product", "A test product", price, 10)
	require.NoError(t, err)

	// Act: Convert entity to item
	item := ProductItemFromEntity(product)

	// Assert: Check item structure
	assert.Equal(t, "PRODUCT#test-product-123", item.PK)
	assert.Equal(t, "PRODUCT#test-product-123", item.SK)
	assert.Equal(t, "PRODUCT#ALL", item.GSI1PK)
	assert.Equal(t, "PRODUCT#test-product-123", item.GSI1SK)
	assert.Equal(t, "PRODUCT", item.Type)
	assert.Equal(t, "test-product-123", item.ID)
	assert.Equal(t, "Test Product", item.Name)
	assert.Equal(t, "A test product", item.Description)
	assert.Equal(t, 1299, item.Price)
	assert.Equal(t, 10, item.Stock)
	assert.WithinDuration(t, time.Now(), item.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), item.UpdatedAt, time.Second)

	// Act: Convert item back to entity
	convertedProduct, err := item.ToEntity()
	require.NoError(t, err)

	// Assert: Check entity structure
	assert.Equal(t, productID, convertedProduct.ID())
	assert.Equal(t, "Test Product", convertedProduct.Name())
	assert.Equal(t, "A test product", convertedProduct.Description())
	assert.Equal(t, price, convertedProduct.Price())
	assert.Equal(t, 10, convertedProduct.Stock())
}

// TestDynamoProductRepository runs integration tests against DynamoDB Local
func TestDynamoProductRepository(t *testing.T) {
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
	repo := NewDynamoProductRepository(client)

	t.Run("save and find product", func(t *testing.T) {
		productID, err := value.NewProductID("test-save-product")
		require.NoError(t, err)

		price, err := value.NewMoney(999)
		require.NoError(t, err)

		product, err := entity.NewProduct(productID, "Save Test Product", "Test product for save operation", price, 5)
		require.NoError(t, err)

		// Save product
		err = repo.Save(ctx, product)
		require.NoError(t, err)

		// Find by ID
		found, err := repo.FindByID(ctx, productID)
		require.NoError(t, err)
		assert.Equal(t, product.ID(), found.ID())
		assert.Equal(t, product.Name(), found.Name())
		assert.Equal(t, product.Description(), found.Description())
		assert.Equal(t, product.Price(), found.Price())
		assert.Equal(t, product.Stock(), found.Stock())

		// Clean up
		err = repo.Delete(ctx, productID)
		assert.NoError(t, err)
	})

	t.Run("find by ID not found", func(t *testing.T) {
		productID, err := value.NewProductID("non-existent-product")
		require.NoError(t, err)

		// Act
		product, err := repo.FindByID(ctx, productID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "product not found")
	})

	t.Run("find all products", func(t *testing.T) {
		// Create test products
		products := make([]*entity.Product, 3)
		for i := 0; i < 3; i++ {
			productID, err := value.NewProductID(fmt.Sprintf("test-product-%d", i))
			require.NoError(t, err)

			price, err := value.NewMoney(int64(1000 + i*100))
			require.NoError(t, err)

			product, err := entity.NewProduct(
				productID,
				fmt.Sprintf("Test Product %d", i),
				fmt.Sprintf("Description for product %d", i),
				price,
				10+i,
			)
			require.NoError(t, err)

			products[i] = product

			// Save each product
			err = repo.Save(ctx, product)
			require.NoError(t, err)
		}

		// Act
		found, err := repo.FindAll(ctx)

		// Assert
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(found), 3) // May have other products from other tests

		// Verify our test products are present
		foundIDs := make([]string, len(found))
		for i, p := range found {
			foundIDs[i] = p.ID().String()
		}

		for _, original := range products {
			assert.Contains(t, foundIDs, original.ID().String())
		}

		// Clean up
		for _, product := range products {
			repo.Delete(ctx, product.ID())
		}
	})

	t.Run("delete product", func(t *testing.T) {
		productID, err := value.NewProductID("test-delete-product")
		require.NoError(t, err)

		price, err := value.NewMoney(799)
		require.NoError(t, err)

		product, err := entity.NewProduct(productID, "Delete Test Product", "Test product for delete operation", price, 3)
		require.NoError(t, err)

		// Save first
		err = repo.Save(ctx, product)
		require.NoError(t, err)

		// Verify it exists
		found, err := repo.FindByID(ctx, productID)
		require.NoError(t, err)
		assert.NotNil(t, found)

		// Act: Delete
		err = repo.Delete(ctx, productID)
		require.NoError(t, err)

		// Assert: Verify it's gone
		deleted, err := repo.FindByID(ctx, productID)
		assert.Error(t, err)
		assert.Nil(t, deleted)
		assert.Contains(t, err.Error(), "product not found")
	})
}
