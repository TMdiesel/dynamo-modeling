package controller

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"dynamo-modeling/internal/adapter/openapi"
	"dynamo-modeling/internal/adapter/presenter"
	"dynamo-modeling/internal/usecase"
)

// ProductController handles product-related requests
type ProductController struct {
	createProductUseCase *usecase.CreateProductUseCase
	getProductUseCase    *usecase.GetProductUseCase
	listProductsUseCase  *usecase.ListProductsUseCase
	updateProductUseCase *usecase.UpdateProductUseCase
	deleteProductUseCase *usecase.DeleteProductUseCase
	presenter            *presenter.ProductPresenter
}

// NewProductController creates a new product controller
func NewProductController(
	createProductUseCase *usecase.CreateProductUseCase,
	getProductUseCase *usecase.GetProductUseCase,
	listProductsUseCase *usecase.ListProductsUseCase,
	updateProductUseCase *usecase.UpdateProductUseCase,
	deleteProductUseCase *usecase.DeleteProductUseCase,
	presenter *presenter.ProductPresenter,
) *ProductController {
	return &ProductController{
		createProductUseCase: createProductUseCase,
		getProductUseCase:    getProductUseCase,
		listProductsUseCase:  listProductsUseCase,
		updateProductUseCase: updateProductUseCase,
		deleteProductUseCase: deleteProductUseCase,
		presenter:            presenter,
	}
}

// CreateProduct handles product creation
func (c *ProductController) CreateProduct(ctx echo.Context) error {
	// 1. リクエスト解析・バリデーション
	var request openapi.ProductRequest
	if err := ctx.Bind(&request); err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Invalid request body")
	}

	// 2. UseCase呼び出し
	command := usecase.CreateProductCommand{
		Name:        request.Name,
		Description: request.Description,
		Price:       int64(request.Price),
		Stock:       request.Stock,
	}

	product, err := c.createProductUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "creation_failed", err.Error())
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentProduct(ctx, http.StatusCreated, product)
}

// GetProduct handles getting a product by ID
func (c *ProductController) GetProduct(ctx echo.Context, productId string) error {
	// 1. バリデーション
	if productId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Product ID is required")
	}

	// 2. UseCase呼び出し
	command := usecase.GetProductCommand{
		ProductID: productId,
	}

	product, err := c.getProductUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusNotFound, "not_found", "Product not found")
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentProduct(ctx, http.StatusOK, product)
}

// ListProducts handles listing all products
func (c *ProductController) ListProducts(ctx echo.Context, params openapi.ListProductsParams) error {
	// 1. パラメータバリデーション
	limit := 100 // デフォルト値
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}

	// 2. UseCase呼び出し
	command := usecase.ListProductsCommand{
		Limit: limit,
	}

	products, err := c.listProductsUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusInternalServerError, "list_failed", err.Error())
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentProducts(ctx, http.StatusOK, products)
}

// UpdateProduct handles product update
func (c *ProductController) UpdateProduct(ctx echo.Context, productId string) error {
	// 1. リクエスト解析・バリデーション
	var request openapi.ProductRequest
	if err := ctx.Bind(&request); err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Invalid request body")
	}

	if productId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Product ID is required")
	}

	// 2. UseCase呼び出し
	command := usecase.UpdateProductCommand{
		ProductID:   productId,
		Name:        request.Name,
		Description: request.Description,
		Price:       int64(request.Price),
		Stock:       request.Stock,
	}

	product, err := c.updateProductUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "update_failed", err.Error())
	}

	// 3. Presenter呼び出し
	return c.presenter.PresentProduct(ctx, http.StatusOK, product)
}

// DeleteProduct handles product deletion
func (c *ProductController) DeleteProduct(ctx echo.Context, productId string) error {
	// 1. バリデーション
	if productId == "" {
		return c.presenter.PresentError(ctx, http.StatusBadRequest, "validation_error", "Product ID is required")
	}

	// 2. UseCase呼び出し
	command := usecase.DeleteProductCommand{
		ProductID: productId,
	}

	err := c.deleteProductUseCase.Execute(context.Background(), command)
	if err != nil {
		return c.presenter.PresentError(ctx, http.StatusInternalServerError, "deletion_failed", err.Error())
	}

	// 3. 成功レスポンス（204 No Content）
	return ctx.NoContent(http.StatusNoContent)
}
