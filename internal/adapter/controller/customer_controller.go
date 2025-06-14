package controller

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"dynamo-modeling/internal/adapter/openapi"
	"dynamo-modeling/internal/adapter/presenter"
	"dynamo-modeling/internal/usecase"
)

// CustomerController handles customer-related requests
type CustomerController struct {
	createCustomerUseCase *usecase.CreateCustomerUseCase
	getCustomerUseCase    *usecase.GetCustomerUseCase
	listCustomersUseCase  *usecase.ListCustomersUseCase
	updateCustomerUseCase *usecase.UpdateCustomerUseCase
	deleteCustomerUseCase *usecase.DeleteCustomerUseCase
	presenter             *presenter.CustomerPresenter
}

// NewCustomerController creates a new customer controller
func NewCustomerController(
	createCustomerUseCase *usecase.CreateCustomerUseCase,
	getCustomerUseCase *usecase.GetCustomerUseCase,
	listCustomersUseCase *usecase.ListCustomersUseCase,
	updateCustomerUseCase *usecase.UpdateCustomerUseCase,
	deleteCustomerUseCase *usecase.DeleteCustomerUseCase,
	presenter *presenter.CustomerPresenter,
) *CustomerController {
	return &CustomerController{
		createCustomerUseCase: createCustomerUseCase,
		getCustomerUseCase:    getCustomerUseCase,
		listCustomersUseCase:  listCustomersUseCase,
		updateCustomerUseCase: updateCustomerUseCase,
		deleteCustomerUseCase: deleteCustomerUseCase,
		presenter:             presenter,
	}
}

// CreateCustomer handles customer creation with validation and business logic
func (c *CustomerController) CreateCustomer(ctx echo.Context) error {
	// 1. リクエスト解析・バリデーション
	var request openapi.CustomerRequest
	if err := ctx.Bind(&request); err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	// バリデーション
	if request.Name == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Name is required")
	}
	if request.Email == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Email is required")
	}

	// 2. UseCase呼び出し
	command := usecase.CreateCustomerCommand{
		Name:  request.Name,
		Email: string(request.Email),
	}

	customer, err := c.createCustomerUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusInternalServerError, "creation_failed", err.Error())
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentCustomer(ctx, http.StatusCreated, customer)
}

// GetCustomer handles getting a customer by ID
func (c *CustomerController) GetCustomer(ctx echo.Context, customerId string) error {
	// 1. バリデーション
	if customerId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Customer ID is required")
	}

	// 2. UseCase呼び出し
	command := usecase.GetCustomerCommand{
		CustomerID: customerId,
	}

	customer, err := c.getCustomerUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusNotFound, "not_found", "Customer not found")
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentCustomer(ctx, http.StatusOK, customer)
}

// ListCustomers handles listing all customers
func (c *CustomerController) ListCustomers(ctx echo.Context, params openapi.ListCustomersParams) error {
	// 1. パラメータバリデーション（必要に応じて）
	limit := 100 // デフォルト値
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}

	// 2. UseCase呼び出し
	command := usecase.ListCustomersCommand{
		Limit: limit,
	}

	customers, err := c.listCustomersUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusInternalServerError, "list_failed", err.Error())
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentCustomers(ctx, http.StatusOK, customers)
}

// UpdateCustomer handles customer update
func (c *CustomerController) UpdateCustomer(ctx echo.Context, customerId string) error {
	// 1. リクエスト解析・バリデーション
	var request openapi.CustomerRequest
	if err := ctx.Bind(&request); err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if customerId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Customer ID is required")
	}

	// 2. UseCase呼び出し
	command := usecase.UpdateCustomerCommand{
		CustomerID: customerId,
		Name:       request.Name,
		Email:      string(request.Email),
	}

	customer, err := c.updateCustomerUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusInternalServerError, "update_failed", err.Error())
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentCustomer(ctx, http.StatusOK, customer)
}

// DeleteCustomer handles customer deletion
func (c *CustomerController) DeleteCustomer(ctx echo.Context, customerId string) error {
	// 1. バリデーション
	if customerId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Customer ID is required")
	}

	// 2. UseCase呼び出し
	command := usecase.DeleteCustomerCommand{
		CustomerID: customerId,
	}

	err := c.deleteCustomerUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusInternalServerError, "deletion_failed", err.Error())
	}

	// 3. レスポンス（204 No Content）
	return ctx.NoContent(http.StatusNoContent)
}
