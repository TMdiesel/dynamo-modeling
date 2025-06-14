package usecase

import (
	"context"

	"dynamo-modeling/internal/domain"
	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/repository"
	"dynamo-modeling/internal/domain/value"
)

// CreateOrderUseCase handles order creation business logic
type CreateOrderUseCase struct {
	orderRepo    repository.OrderRepository
	customerRepo repository.CustomerRepository
	productRepo  repository.ProductRepository
}

// CreateOrderCommand represents the input for creating an order
type CreateOrderCommand struct {
	CustomerID string
	Items      []CreateOrderItemCommand
}

// CreateOrderItemCommand represents an order item
type CreateOrderItemCommand struct {
	ProductID string
	Quantity  int
}

// NewCreateOrderUseCase creates a new create order use case
func NewCreateOrderUseCase(
	orderRepo repository.OrderRepository,
	customerRepo repository.CustomerRepository,
	productRepo repository.ProductRepository,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		productRepo:  productRepo,
	}
}

// Execute executes the create order use case
func (uc *CreateOrderUseCase) Execute(ctx context.Context, cmd CreateOrderCommand) (*entity.Order, error) {
	// 1. 顧客の存在確認
	customerID := value.CustomerID(cmd.CustomerID)
	customerExists, err := uc.customerRepo.Exists(ctx, customerID)
	if err != nil {
		return nil, domain.RepositoryError("failed to check customer existence", err)
	}
	if !customerExists {
		return nil, domain.NewDomainError("CUSTOMER_NOT_FOUND", "Customer not found", nil)
	}

	// 2. 注文商品の検証と在庫確認
	var orderItems []entity.OrderItem
	for _, itemCmd := range cmd.Items {
		productID := value.ProductID(itemCmd.ProductID)

		// 商品の存在確認
		product, err := uc.productRepo.FindByID(ctx, productID)
		if err != nil {
			return nil, domain.RepositoryError("failed to find product", err)
		}
		if product == nil {
			return nil, domain.NewDomainError("PRODUCT_NOT_FOUND", "Product not found: "+itemCmd.ProductID, nil)
		}

		// 在庫確認
		if !product.IsInStock(itemCmd.Quantity) {
			return nil, domain.NewDomainError("INSUFFICIENT_STOCK",
				"Insufficient stock for product: "+itemCmd.ProductID, nil)
		}

		// 注文アイテム作成
		orderItem, err := entity.NewOrderItem(productID, itemCmd.Quantity, product.Price())
		if err != nil {
			return nil, domain.InvalidInputError("invalid order item: " + err.Error())
		}

		orderItems = append(orderItems, *orderItem)
	}

	// 3. 新しいOrder IDを生成
	orderID := value.GenerateOrderID()

	// 4. 注文エンティティ作成
	order, err := entity.NewOrder(orderID, customerID, orderItems)
	if err != nil {
		return nil, domain.InvalidInputError("failed to create order: " + err.Error())
	}

	// 5. 在庫の予約（商品の在庫を減らす）
	for _, itemCmd := range cmd.Items {
		productID := value.ProductID(itemCmd.ProductID)
		product, _ := uc.productRepo.FindByID(ctx, productID)

		err = product.ReserveStock(itemCmd.Quantity)
		if err != nil {
			return nil, domain.NewDomainError("STOCK_RESERVATION_FAILED",
				"Failed to reserve stock: "+err.Error(), nil)
		}

		// 商品の在庫更新をDBに反映
		err = uc.productRepo.Save(ctx, product)
		if err != nil {
			return nil, domain.RepositoryError("failed to update product stock", err)
		}
	}

	// 6. 注文をリポジトリに保存
	err = uc.orderRepo.Save(ctx, order)
	if err != nil {
		return nil, domain.RepositoryError("failed to save order", err)
	}

	return order, nil
}

// GetOrderUseCase handles getting an order by ID
type GetOrderUseCase struct {
	orderRepo repository.OrderRepository
}

// GetOrderCommand represents the input for getting an order
type GetOrderCommand struct {
	OrderID string
}

// NewGetOrderUseCase creates a new get order use case
func NewGetOrderUseCase(orderRepo repository.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{
		orderRepo: orderRepo,
	}
}

// Execute executes the get order use case
func (uc *GetOrderUseCase) Execute(ctx context.Context, cmd GetOrderCommand) (*entity.Order, error) {
	// 1. 値オブジェクトの作成・バリデーション
	orderID := value.OrderID(cmd.OrderID)

	// 2. リポジトリから取得
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, domain.NewDomainError("ORDER_NOT_FOUND", "Order not found", nil)
	}

	return order, nil
}

// ListOrdersUseCase handles listing orders
type ListOrdersUseCase struct {
	orderRepo repository.OrderRepository
}

// ListOrdersCommand represents the input for listing orders
type ListOrdersCommand struct {
	CustomerID *string
	Limit      int
}

// NewListOrdersUseCase creates a new list orders use case
func NewListOrdersUseCase(orderRepo repository.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		orderRepo: orderRepo,
	}
}

// Execute executes the list orders use case
func (uc *ListOrdersUseCase) Execute(ctx context.Context, cmd ListOrdersCommand) ([]*entity.Order, error) {
	// ビジネスロジック: デフォルトの制限値設定
	if cmd.Limit <= 0 || cmd.Limit > 1000 {
		cmd.Limit = 100 // デフォルト制限
	}

	// 顧客別 or 全体の注文リスト取得
	if cmd.CustomerID != nil && *cmd.CustomerID != "" {
		customerID := value.CustomerID(*cmd.CustomerID)
		orders, _, err := uc.orderRepo.FindByCustomerID(ctx, customerID, cmd.Limit, nil)
		if err != nil {
			return nil, domain.RepositoryError("failed to list orders by customer", err)
		}
		return orders, nil
	}

	// 全注文リスト（将来的にはページネーション対応予定）
	// 現在は簡単な実装として空配列を返す
	// TODO: OrderRepository にFindAllメソッドを追加
	return []*entity.Order{}, nil
}

// UpdateOrderStatusUseCase handles order status update
type UpdateOrderStatusUseCase struct {
	orderRepo repository.OrderRepository
}

// UpdateOrderStatusCommand represents the input for updating order status
type UpdateOrderStatusCommand struct {
	OrderID string
	Status  string
}

// NewUpdateOrderStatusUseCase creates a new update order status use case
func NewUpdateOrderStatusUseCase(orderRepo repository.OrderRepository) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{
		orderRepo: orderRepo,
	}
}

// Execute executes the update order status use case
func (uc *UpdateOrderStatusUseCase) Execute(ctx context.Context, cmd UpdateOrderStatusCommand) (*entity.Order, error) {
	// 1. 値オブジェクトの作成・バリデーション
	orderID := value.OrderID(cmd.OrderID)

	// 2. 既存の注文を取得
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, domain.NewDomainError("ORDER_NOT_FOUND", "Order not found", nil)
	}

	// 3. ステータス更新（ビジネスロジックチェック含む）
	err = order.UpdateStatus(entity.OrderStatus(cmd.Status))
	if err != nil {
		return nil, domain.InvalidInputError("invalid status transition: " + err.Error())
	}

	// 4. リポジトリに保存
	err = uc.orderRepo.Save(ctx, order)
	if err != nil {
		return nil, domain.RepositoryError("failed to update order", err)
	}

	return order, nil
}
