package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/guregu/dynamo/v2"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
	"dynamo-modeling/internal/infrastructure"
)

// DynamoOrderRepository implements OrderRepository using DynamoDB
type DynamoOrderRepository struct {
	client *infrastructure.DynamoDBClient
}

// NewDynamoOrderRepository creates a new DynamoDB order repository
func NewDynamoOrderRepository(client *infrastructure.DynamoDBClient) *DynamoOrderRepository {
	return &DynamoOrderRepository{
		client: client,
	}
}

// OrderItemData represents order item data for DynamoDB storage
type OrderItemData struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unitPrice"` // Price in cents
}

// OrderItem represents an order item in DynamoDB
type OrderItem struct {
	PK         string    `dynamo:"PK"`         // ORDER#{OrderID}
	SK         string    `dynamo:"SK"`         // ORDER#{OrderID}
	GSI1PK     string    `dynamo:"GSI1PK"`     // CUSTOMER#{CustomerID}
	GSI1SK     string    `dynamo:"GSI1SK"`     // ORDER#{CreatedAt}#{OrderID}
	Type       string    `dynamo:"Type"`       // "ORDER"
	ID         string    `dynamo:"ID"`         // OrderID
	CustomerID string    `dynamo:"CustomerID"` // CustomerID
	Items      string    `dynamo:"Items"`      // JSON array of OrderItemData
	Status     string    `dynamo:"Status"`     // Order status
	Total      int64     `dynamo:"Total"`      // Total price in cents
	CreatedAt  time.Time `dynamo:"CreatedAt"`  // Creation timestamp
	UpdatedAt  time.Time `dynamo:"UpdatedAt"`  // Last update timestamp
}

// ToEntity converts OrderItem to Order entity
func (item *OrderItem) ToEntity() (*entity.Order, error) {
	orderID, err := value.NewOrderID(item.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	customerID, err := value.NewCustomerID(item.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	// Parse order items
	var itemsData []OrderItemData
	if err := json.Unmarshal([]byte(item.Items), &itemsData); err != nil {
		return nil, fmt.Errorf("failed to parse order items: %w", err)
	}

	orderItems := make([]entity.OrderItem, 0, len(itemsData))
	for _, data := range itemsData {
		productID, err := value.NewProductID(data.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID in order item: %w", err)
		}

		unitPrice, err := value.NewMoney(data.UnitPrice)
		if err != nil {
			return nil, fmt.Errorf("invalid unit price in order item: %w", err)
		}

		orderItem, err := entity.NewOrderItem(productID, data.Quantity, unitPrice)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}

		orderItems = append(orderItems, *orderItem)
	}

	order, err := entity.NewOrder(orderID, customerID, orderItems)
	if err != nil {
		return nil, fmt.Errorf("failed to create order entity: %w", err)
	}

	return order, nil
}

// FromEntity converts Order entity to OrderItem
func OrderItemFromEntity(order *entity.Order) (*OrderItem, error) {
	orderID := order.ID().String()
	customerID := order.CustomerID().String()

	// Convert order items to JSON
	itemsData := make([]OrderItemData, 0, len(order.Items()))
	for _, item := range order.Items() {
		itemsData = append(itemsData, OrderItemData{
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice.Cents(),
		})
	}

	itemsJSON, err := json.Marshal(itemsData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order items: %w", err)
	}

	return &OrderItem{
		PK:         fmt.Sprintf("ORDER#%s", orderID),
		SK:         fmt.Sprintf("ORDER#%s", orderID),
		GSI1PK:     fmt.Sprintf("CUSTOMER#%s", customerID),
		GSI1SK:     fmt.Sprintf("ORDER#%s#%s", order.CreatedAt().Format(time.RFC3339), orderID),
		Type:       "ORDER",
		ID:         orderID,
		CustomerID: customerID,
		Items:      string(itemsJSON),
		Status:     string(order.Status()),
		Total:      order.Total().Cents(),
		CreatedAt:  order.CreatedAt(),
		UpdatedAt:  order.UpdatedAt(),
	}, nil
}

// Save creates or updates an order
func (r *DynamoOrderRepository) Save(ctx context.Context, order *entity.Order) error {
	slog.Info("Saving order", "orderID", order.ID().String())

	item, err := OrderItemFromEntity(order)
	if err != nil {
		return fmt.Errorf("failed to convert order to item: %w", err)
	}

	table := r.client.GetTable()

	err = table.Put(item).Run(ctx)
	if err != nil {
		slog.Error("Failed to save order", "orderID", order.ID().String(), "error", err)
		return fmt.Errorf("failed to save order: %w", err)
	}

	slog.Info("Order saved successfully", "orderID", order.ID().String())
	return nil
}

// FindByID retrieves an order by ID
func (r *DynamoOrderRepository) FindByID(ctx context.Context, id value.OrderID) (*entity.Order, error) {
	slog.Info("Finding order by ID", "orderID", id.String())

	var item OrderItem
	table := r.client.GetTable()

	err := table.Get("PK", fmt.Sprintf("ORDER#%s", id.String())).
		Range("SK", dynamo.Equal, fmt.Sprintf("ORDER#%s", id.String())).
		One(ctx, &item)

	if err != nil {
		if err == dynamo.ErrNotFound {
			slog.Info("Order not found", "orderID", id.String())
			return nil, fmt.Errorf("order not found: %s", id.String())
		}
		slog.Error("Failed to find order", "orderID", id.String(), "error", err)
		return nil, fmt.Errorf("failed to find order: %w", err)
	}

	order, err := item.ToEntity()
	if err != nil {
		slog.Error("Failed to convert item to entity", "orderID", id.String(), "error", err)
		return nil, fmt.Errorf("failed to convert item to entity: %w", err)
	}

	slog.Info("Order found successfully", "orderID", id.String())
	return order, nil
}

// FindByCustomerID retrieves all orders for a customer
func (r *DynamoOrderRepository) FindByCustomerID(ctx context.Context, customerID value.CustomerID) ([]*entity.Order, error) {
	slog.Info("Finding orders by customer ID", "customerID", customerID.String())

	var items []OrderItem
	table := r.client.GetTable()

	err := table.Get("GSI1PK", fmt.Sprintf("CUSTOMER#%s", customerID.String())).
		Index("GSI1").
		All(ctx, &items)

	if err != nil {
		slog.Error("Failed to find orders by customer ID", "customerID", customerID.String(), "error", err)
		return nil, fmt.Errorf("failed to find orders by customer ID: %w", err)
	}

	orders := make([]*entity.Order, 0, len(items))
	for _, item := range items {
		order, err := item.ToEntity()
		if err != nil {
			slog.Error("Failed to convert item to entity", "orderID", item.ID, "error", err)
			continue // Skip invalid items
		}
		orders = append(orders, order)
	}

	slog.Info("Found orders successfully", "customerID", customerID.String(), "count", len(orders))
	return orders, nil
}

// Delete removes an order
func (r *DynamoOrderRepository) Delete(ctx context.Context, id value.OrderID) error {
	slog.Info("Deleting order", "orderID", id.String())

	table := r.client.GetTable()

	err := table.Delete("PK", fmt.Sprintf("ORDER#%s", id.String())).
		Range("SK", fmt.Sprintf("ORDER#%s", id.String())).
		Run(ctx)

	if err != nil {
		slog.Error("Failed to delete order", "orderID", id.String(), "error", err)
		return fmt.Errorf("failed to delete order: %w", err)
	}

	slog.Info("Order deleted successfully", "orderID", id.String())
	return nil
}
