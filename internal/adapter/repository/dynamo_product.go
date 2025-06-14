package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/guregu/dynamo/v2"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
	"dynamo-modeling/internal/infrastructure"
)

// DynamoProductRepository implements ProductRepository using DynamoDB
type DynamoProductRepository struct {
	client *infrastructure.DynamoDBClient
}

// NewDynamoProductRepository creates a new DynamoDB product repository
func NewDynamoProductRepository(client *infrastructure.DynamoDBClient) *DynamoProductRepository {
	return &DynamoProductRepository{
		client: client,
	}
}

// ProductItem represents a product item in DynamoDB
type ProductItem struct {
	PK          string    `dynamo:"PK"`          // PRODUCT#{ProductID}
	SK          string    `dynamo:"SK"`          // PRODUCT#{ProductID}
	GSI1PK      string    `dynamo:"GSI1PK"`      // CATEGORY#{CategoryName} (for future use)
	GSI1SK      string    `dynamo:"GSI1SK"`      // PRODUCT#{ProductID}
	Type        string    `dynamo:"Type"`        // "PRODUCT"
	ID          string    `dynamo:"ID"`          // ProductID
	Name        string    `dynamo:"Name"`        // Product name
	Description string    `dynamo:"Description"` // Product description
	Price       int       `dynamo:"Price"`       // Price in cents (Money value)
	Stock       int       `dynamo:"Stock"`       // Stock quantity
	CreatedAt   time.Time `dynamo:"CreatedAt"`   // Creation timestamp
	UpdatedAt   time.Time `dynamo:"UpdatedAt"`   // Last update timestamp
}

// ToEntity converts ProductItem to Product entity
func (item *ProductItem) ToEntity() (*entity.Product, error) {
	productID, err := value.NewProductID(item.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	price, err := value.NewMoney(int64(item.Price))
	if err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}

	product, err := entity.NewProduct(productID, item.Name, item.Description, price, item.Stock)
	if err != nil {
		return nil, fmt.Errorf("failed to create product entity: %w", err)
	}

	return product, nil
}

// FromEntity converts Product entity to ProductItem
func ProductItemFromEntity(product *entity.Product) *ProductItem {
	productID := product.ID().String()

	return &ProductItem{
		PK:          fmt.Sprintf("PRODUCT#%s", productID),
		SK:          fmt.Sprintf("PRODUCT#%s", productID),
		GSI1PK:      "PRODUCT#ALL", // For listing all products
		GSI1SK:      fmt.Sprintf("PRODUCT#%s", productID),
		Type:        "PRODUCT",
		ID:          productID,
		Name:        product.Name(),
		Description: product.Description(),
		Price:       int(product.Price().Cents()),
		Stock:       product.Stock(),
		CreatedAt:   product.CreatedAt(),
		UpdatedAt:   product.UpdatedAt(),
	}
}

// Save creates or updates a product
func (r *DynamoProductRepository) Save(ctx context.Context, product *entity.Product) error {
	slog.Info("Saving product", "productID", product.ID().String())

	item := ProductItemFromEntity(product)
	table := r.client.GetTable()

	err := table.Put(item).Run(ctx)
	if err != nil {
		slog.Error("Failed to save product", "productID", product.ID().String(), "error", err)
		return fmt.Errorf("failed to save product: %w", err)
	}

	slog.Info("Product saved successfully", "productID", product.ID().String())
	return nil
}

// FindByID retrieves a product by their ID
func (r *DynamoProductRepository) FindByID(ctx context.Context, id value.ProductID) (*entity.Product, error) {
	slog.Info("Finding product by ID", "productID", id.String())

	var item ProductItem
	table := r.client.GetTable()

	err := table.Get("PK", fmt.Sprintf("PRODUCT#%s", id.String())).
		Range("SK", dynamo.Equal, fmt.Sprintf("PRODUCT#%s", id.String())).
		One(ctx, &item)

	if err != nil {
		if err == dynamo.ErrNotFound {
			slog.Info("Product not found", "productID", id.String())
			return nil, fmt.Errorf("product not found: %s", id.String())
		}
		slog.Error("Failed to find product", "productID", id.String(), "error", err)
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	product, err := item.ToEntity()
	if err != nil {
		slog.Error("Failed to convert item to entity", "productID", id.String(), "error", err)
		return nil, fmt.Errorf("failed to convert item to entity: %w", err)
	}

	slog.Info("Product found successfully", "productID", id.String())
	return product, nil
}

// FindAll retrieves all products
func (r *DynamoProductRepository) FindAll(ctx context.Context) ([]*entity.Product, error) {
	slog.Info("Finding all products")

	var items []ProductItem
	table := r.client.GetTable()

	err := table.Get("GSI1PK", "PRODUCT#ALL").
		Index("GSI1").
		All(ctx, &items)

	if err != nil {
		slog.Error("Failed to find all products", "error", err)
		return nil, fmt.Errorf("failed to find all products: %w", err)
	}

	products := make([]*entity.Product, 0, len(items))
	for _, item := range items {
		product, err := item.ToEntity()
		if err != nil {
			slog.Error("Failed to convert item to entity", "productID", item.ID, "error", err)
			continue // Skip invalid items
		}
		products = append(products, product)
	}

	slog.Info("Found products successfully", "count", len(products))
	return products, nil
}

// Delete removes a product
func (r *DynamoProductRepository) Delete(ctx context.Context, id value.ProductID) error {
	slog.Info("Deleting product", "productID", id.String())

	table := r.client.GetTable()

	err := table.Delete("PK", fmt.Sprintf("PRODUCT#%s", id.String())).
		Range("SK", fmt.Sprintf("PRODUCT#%s", id.String())).
		Run(ctx)

	if err != nil {
		slog.Error("Failed to delete product", "productID", id.String(), "error", err)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	slog.Info("Product deleted successfully", "productID", id.String())
	return nil
}
