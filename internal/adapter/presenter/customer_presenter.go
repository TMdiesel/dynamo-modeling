package presenter

import (
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"dynamo-modeling/internal/adapter/openapi"
	"dynamo-modeling/internal/domain/entity"
)

// CustomerPresenter handles customer response presentation
type CustomerPresenter struct{}

// NewCustomerPresenter creates a new customer presenter
func NewCustomerPresenter() *CustomerPresenter {
	return &CustomerPresenter{}
}

// PresentCustomer presents a single customer
func (p *CustomerPresenter) PresentCustomer(ctx echo.Context, statusCode int, customer *entity.Customer) error {
	response := openapi.CustomerResponse{
		Id:        customer.ID().String(),
		Name:      customer.Name(),
		Email:     openapi_types.Email(customer.Email().String()),
		CreatedAt: customer.CreatedAt(),
		UpdatedAt: customer.UpdatedAt(),
	}

	return ctx.JSON(statusCode, response)
}

// PresentCustomers presents a list of customers
func (p *CustomerPresenter) PresentCustomers(ctx echo.Context, statusCode int, customers []*entity.Customer) error {
	responses := make([]openapi.CustomerResponse, len(customers))

	for i, customer := range customers {
		responses[i] = openapi.CustomerResponse{
			Id:        customer.ID().String(),
			Name:      customer.Name(),
			Email:     openapi_types.Email(customer.Email().String()),
			CreatedAt: customer.CreatedAt(),
			UpdatedAt: customer.UpdatedAt(),
		}
	}

	return ctx.JSON(statusCode, responses)
}

// PresentError presents an error response
func (p *CustomerPresenter) PresentError(ctx echo.Context, statusCode int, code, message string) error {
	errorResponse := openapi.Error{
		Code:    code,
		Message: message,
	}

	return ctx.JSON(statusCode, errorResponse)
}
