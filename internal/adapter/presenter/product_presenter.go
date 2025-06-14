package presenter

import (
	"github.com/labstack/echo/v4"

	"dynamo-modeling/internal/adapter/openapi"
	"dynamo-modeling/internal/domain/entity"
)

// ProductPresenter handles product response presentation
type ProductPresenter struct{}

// NewProductPresenter creates a new product presenter
func NewProductPresenter() *ProductPresenter {
	return &ProductPresenter{}
}

// PresentProduct presents a single product
func (p *ProductPresenter) PresentProduct(ctx echo.Context, statusCode int, product *entity.Product) error {
	response := openapi.ProductResponse{
		Id:          product.ID().String(),
		Name:        product.Name(),
		Description: product.Description(),
		Price:       int(product.Price().Cents()),
		Stock:       product.Stock(),
		CreatedAt:   product.CreatedAt(),
		UpdatedAt:   product.UpdatedAt(),
	}

	return ctx.JSON(statusCode, response)
}

// PresentProducts presents a list of products
func (p *ProductPresenter) PresentProducts(ctx echo.Context, statusCode int, products []*entity.Product) error {
	responses := make([]openapi.ProductResponse, len(products))

	for i, product := range products {
		responses[i] = openapi.ProductResponse{
			Id:          product.ID().String(),
			Name:        product.Name(),
			Description: product.Description(),
			Price:       int(product.Price().Cents()),
			Stock:       product.Stock(),
			CreatedAt:   product.CreatedAt(),
			UpdatedAt:   product.UpdatedAt(),
		}
	}

	return ctx.JSON(statusCode, responses)
}

// PresentError presents an error response
func (p *ProductPresenter) PresentError(ctx echo.Context, statusCode int, code, message string) error {
	errorResponse := openapi.Error{
		Code:    code,
		Message: message,
	}

	return ctx.JSON(statusCode, errorResponse)
}
