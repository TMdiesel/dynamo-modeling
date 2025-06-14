package entity

import (
	"fmt"
	"time"

	"dynamo-modeling/internal/domain/value"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID value.ProductID
	Quantity  int
	UnitPrice value.Money
}

// NewOrderItem creates a new OrderItem
func NewOrderItem(productID value.ProductID, quantity int, unitPrice value.Money) (*OrderItem, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}
	if unitPrice.IsZero() {
		return nil, fmt.Errorf("unit price must be positive")
	}

	return &OrderItem{
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}, nil
}

// TotalPrice calculates the total price for this order item
func (oi *OrderItem) TotalPrice() value.Money {
	total, _ := oi.UnitPrice.Multiply(float64(oi.Quantity))
	return total
}

// Order represents an order entity
type Order struct {
	id         value.OrderID
	customerID value.CustomerID
	items      []OrderItem
	status     OrderStatus
	total      value.Money
	createdAt  time.Time
	updatedAt  time.Time
}

// NewOrder creates a new Order entity
func NewOrder(id value.OrderID, customerID value.CustomerID, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	now := time.Now()
	order := &Order{
		id:         id,
		customerID: customerID,
		items:      make([]OrderItem, len(items)),
		status:     OrderStatusPending,
		createdAt:  now,
		updatedAt:  now,
	}

	// Copy items and calculate total
	copy(order.items, items)
	order.calculateTotal()

	return order, nil
}

// ID returns the order ID
func (o *Order) ID() value.OrderID {
	return o.id
}

// CustomerID returns the customer ID
func (o *Order) CustomerID() value.CustomerID {
	return o.customerID
}

// Items returns a copy of the order items
func (o *Order) Items() []OrderItem {
	items := make([]OrderItem, len(o.items))
	copy(items, o.items)
	return items
}

// Status returns the order status
func (o *Order) Status() OrderStatus {
	return o.status
}

// Total returns the order total
func (o *Order) Total() value.Money {
	return o.total
}

// CreatedAt returns the creation timestamp
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// UpdatedAt returns the last update timestamp
func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

// Confirm changes the order status to confirmed
func (o *Order) Confirm() error {
	if o.status != OrderStatusPending {
		return fmt.Errorf("can only confirm pending orders, current status: %s", o.status)
	}
	o.status = OrderStatusConfirmed
	o.updatedAt = time.Now()
	return nil
}

// Ship changes the order status to shipped
func (o *Order) Ship() error {
	if o.status != OrderStatusConfirmed {
		return fmt.Errorf("can only ship confirmed orders, current status: %s", o.status)
	}
	o.status = OrderStatusShipped
	o.updatedAt = time.Now()
	return nil
}

// Deliver changes the order status to delivered
func (o *Order) Deliver() error {
	if o.status != OrderStatusShipped {
		return fmt.Errorf("can only deliver shipped orders, current status: %s", o.status)
	}
	o.status = OrderStatusDelivered
	o.updatedAt = time.Now()
	return nil
}

// Cancel changes the order status to cancelled
func (o *Order) Cancel() error {
	if o.status == OrderStatusDelivered || o.status == OrderStatusCancelled {
		return fmt.Errorf("cannot cancel order with status: %s", o.status)
	}
	o.status = OrderStatusCancelled
	o.updatedAt = time.Now()
	return nil
}

// IsPending checks if the order is pending
func (o *Order) IsPending() bool {
	return o.status == OrderStatusPending
}

// IsConfirmed checks if the order is confirmed
func (o *Order) IsConfirmed() bool {
	return o.status == OrderStatusConfirmed
}

// IsCancelled checks if the order is cancelled
func (o *Order) IsCancelled() bool {
	return o.status == OrderStatusCancelled
}

// ItemCount returns the total number of items in the order
func (o *Order) ItemCount() int {
	total := 0
	for _, item := range o.items {
		total += item.Quantity
	}
	return total
}

// calculateTotal calculates the total price of the order
func (o *Order) calculateTotal() {
	total, _ := value.NewMoney(0)
	for _, item := range o.items {
		total = total.Add(item.TotalPrice())
	}
	o.total = total
}

// Equals compares two Order entities by ID
func (o *Order) Equals(other *Order) bool {
	if other == nil {
		return false
	}
	return o.id == other.id
}
