package handler

import (
	"github.com/labstack/echo/v4"

	"dynamo-modeling/internal/adapter/controller"
	"dynamo-modeling/internal/adapter/openapi"
)

// APIHandler implements the OpenAPI server interface
type APIHandler struct {
	customerController *controller.CustomerController
	productController  *controller.ProductController
	orderController    *controller.OrderController
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(
	customerController *controller.CustomerController,
	productController *controller.ProductController,
	orderController *controller.OrderController,
) *APIHandler {
	return &APIHandler{
		customerController: customerController,
		productController:  productController,
		orderController:    orderController,
	}
}

// Verify APIHandler implements ServerInterface
var _ openapi.ServerInterface = (*APIHandler)(nil)

// Customer endpoints

// ListCustomers handles listing all customers
func (h *APIHandler) ListCustomers(ctx echo.Context, params openapi.ListCustomersParams) error {
	return h.customerController.ListCustomers(ctx, params)
}

// CreateCustomer handles customer creation
func (h *APIHandler) CreateCustomer(ctx echo.Context) error {
	return h.customerController.CreateCustomer(ctx)
}

// DeleteCustomer handles customer deletion
func (h *APIHandler) DeleteCustomer(ctx echo.Context, customerId string) error {
	return h.customerController.DeleteCustomer(ctx, customerId)
}

// GetCustomer handles getting a customer by ID
func (h *APIHandler) GetCustomer(ctx echo.Context, customerId string) error {
	return h.customerController.GetCustomer(ctx, customerId)
}

// UpdateCustomer handles customer update
func (h *APIHandler) UpdateCustomer(ctx echo.Context, customerId string) error {
	return h.customerController.UpdateCustomer(ctx, customerId)
}

// GetCustomerOrders handles getting customer orders
func (h *APIHandler) GetCustomerOrders(ctx echo.Context, customerId string, params openapi.GetCustomerOrdersParams) error {
	return h.orderController.GetCustomerOrders(ctx, customerId, params)
}

// Order endpoints

// ListOrders handles listing all orders
func (h *APIHandler) ListOrders(ctx echo.Context, params openapi.ListOrdersParams) error {
	return h.orderController.ListOrders(ctx, params)
}

// CreateOrder handles order creation
func (h *APIHandler) CreateOrder(ctx echo.Context) error {
	return h.orderController.CreateOrder(ctx)
}

// GetOrder handles getting an order by ID
func (h *APIHandler) GetOrder(ctx echo.Context, orderId string) error {
	return h.orderController.GetOrder(ctx, orderId)
}

// UpdateOrderStatus handles order status update
func (h *APIHandler) UpdateOrderStatus(ctx echo.Context, orderId string) error {
	return h.orderController.UpdateOrderStatus(ctx, orderId)
}

// Product endpoints

// ListProducts handles listing all products
func (h *APIHandler) ListProducts(ctx echo.Context, params openapi.ListProductsParams) error {
	return h.productController.ListProducts(ctx, params)
}

// CreateProduct handles product creation
func (h *APIHandler) CreateProduct(ctx echo.Context) error {
	return h.productController.CreateProduct(ctx)
}

// DeleteProduct handles product deletion
func (h *APIHandler) DeleteProduct(ctx echo.Context, productId string) error {
	return h.productController.DeleteProduct(ctx, productId)
}

// GetProduct handles getting a product by ID
func (h *APIHandler) GetProduct(ctx echo.Context, productId string) error {
	return h.productController.GetProduct(ctx, productId)
}

// UpdateProduct handles product update
func (h *APIHandler) UpdateProduct(ctx echo.Context, productId string) error {
	return h.productController.UpdateProduct(ctx, productId)
}
