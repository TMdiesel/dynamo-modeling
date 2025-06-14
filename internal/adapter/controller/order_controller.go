package controller

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"dynamo-modeling/internal/adapter/openapi"
	"dynamo-modeling/internal/adapter/presenter"
	"dynamo-modeling/internal/usecase"
)

// OrderController handles order-related requests
type OrderController struct {
	createOrderUseCase       *usecase.CreateOrderUseCase
	getOrderUseCase          *usecase.GetOrderUseCase
	listOrdersUseCase        *usecase.ListOrdersUseCase
	updateOrderStatusUseCase *usecase.UpdateOrderStatusUseCase
	presenter                *presenter.OrderPresenter
}

// NewOrderController creates a new order controller
func NewOrderController(
	createOrderUseCase *usecase.CreateOrderUseCase,
	getOrderUseCase *usecase.GetOrderUseCase,
	listOrdersUseCase *usecase.ListOrdersUseCase,
	updateOrderStatusUseCase *usecase.UpdateOrderStatusUseCase,
	presenter *presenter.OrderPresenter,
) *OrderController {
	return &OrderController{
		createOrderUseCase:       createOrderUseCase,
		getOrderUseCase:          getOrderUseCase,
		listOrdersUseCase:        listOrdersUseCase,
		updateOrderStatusUseCase: updateOrderStatusUseCase,
		presenter:                presenter,
	}
}

// CreateOrder handles order creation
func (c *OrderController) CreateOrder(ctx echo.Context) error {
	// 1. リクエスト解析・バリデーション
	var request openapi.OrderRequest
	if err := ctx.Bind(&request); err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Invalid request body")
	}

	// 2. コマンド構築
	var items []usecase.CreateOrderItemCommand
	for _, item := range request.Items {
		items = append(items, usecase.CreateOrderItemCommand{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	command := usecase.CreateOrderCommand{
		CustomerID: request.CustomerId,
		Items:      items,
	}

	// 3. UseCase呼び出し
	order, err := c.createOrderUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "creation_failed", err.Error())
	}

	// 4. Presenter呼び出し
	return c.presenter.PresentOrder(ctx, http.StatusCreated, order)
}

// GetOrder handles getting an order by ID
func (c *OrderController) GetOrder(ctx echo.Context, orderId string) error {
	// 1. バリデーション
	if orderId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Order ID is required")
	}

	// 2. UseCase呼び出し
	command := usecase.GetOrderCommand{
		OrderID: orderId,
	}

	order, err := c.getOrderUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusNotFound, "not_found", "Order not found")
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentOrder(ctx, http.StatusOK, order)
}

// ListOrders handles listing orders
func (c *OrderController) ListOrders(ctx echo.Context, params openapi.ListOrdersParams) error {
	// 1. パラメータ処理
	limit := 100 // デフォルト値
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}

	var customerID *string
	if params.CustomerId != nil && *params.CustomerId != "" {
		customerID = params.CustomerId
	}

	// 2. UseCase呼び出し
	command := usecase.ListOrdersCommand{
		CustomerID: customerID,
		Limit:      limit,
	}

	orders, err := c.listOrdersUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusInternalServerError, "list_failed", err.Error())
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentOrders(ctx, http.StatusOK, orders)
}

// GetCustomerOrders handles getting orders by customer ID
func (c *OrderController) GetCustomerOrders(ctx echo.Context, customerId string, params openapi.GetCustomerOrdersParams) error {
	// 1. バリデーション
	if customerId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Customer ID is required")
	}

	// 2. パラメータ処理
	limit := 100 // デフォルト値
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}

	// 3. UseCase呼び出し
	command := usecase.ListOrdersCommand{
		CustomerID: &customerId,
		Limit:      limit,
	}

	orders, err := c.listOrdersUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusInternalServerError, "list_failed", err.Error())
	}

	// 4. Presenter呼び出し
	return c.presenter.PresentOrders(ctx, http.StatusOK, orders)
}

// UpdateOrderStatus handles order status update
func (c *OrderController) UpdateOrderStatus(ctx echo.Context, orderId string) error {
	// 1. リクエスト解析・バリデーション
	var request openapi.UpdateOrderStatusJSONBody
	if err := ctx.Bind(&request); err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Invalid request body")
	}

	if orderId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Order ID is required")
	}

	// 2. UseCase呼び出し
	command := usecase.UpdateOrderStatusCommand{
		OrderID: orderId,
		Status:  string(request.Status),
	}

	order, err := c.updateOrderStatusUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "update_failed", err.Error())
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentOrder(ctx, http.StatusOK, order)
}
