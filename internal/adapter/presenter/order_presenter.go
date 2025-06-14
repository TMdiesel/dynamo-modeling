package presenter

import (
	"github.com/labstack/echo/v4"

	"dynamo-modeling/internal/adapter/openapi"
	"dynamo-modeling/internal/domain/entity"
)

// OrderPresenter handles order response presentation
type OrderPresenter struct{}

// NewOrderPresenter creates a new order presenter
func NewOrderPresenter() *OrderPresenter {
	return &OrderPresenter{}
}

// PresentOrder presents a single order
func (p *OrderPresenter) PresentOrder(ctx echo.Context, statusCode int, order *entity.Order) error {
	// OrderItemsをOpenAPI形式に変換
	items := make([]openapi.OrderItemResponse, len(order.Items()))
	for i, item := range order.Items() {
		items[i] = openapi.OrderItemResponse{
			ProductId:  item.ProductID.String(),
			Quantity:   item.Quantity,
			UnitPrice:  int(item.UnitPrice.Cents()),
			TotalPrice: int(item.UnitPrice.Cents()) * item.Quantity,
		}
	}

	response := openapi.OrderResponse{
		Id:          order.ID().String(),
		CustomerId:  order.CustomerID().String(),
		Items:       items,
		Status:      openapi.OrderResponseStatus(order.Status()),
		TotalAmount: int(order.Total().Cents()),
		CreatedAt:   order.CreatedAt(),
		UpdatedAt:   order.UpdatedAt(),
	}

	return ctx.JSON(statusCode, response)
}

// PresentOrders presents a list of orders
func (p *OrderPresenter) PresentOrders(ctx echo.Context, statusCode int, orders []*entity.Order) error {
	responses := make([]openapi.OrderResponse, len(orders))

	for i, order := range orders {
		// OrderItemsをOpenAPI形式に変換
		items := make([]openapi.OrderItemResponse, len(order.Items()))
		for j, item := range order.Items() {
			items[j] = openapi.OrderItemResponse{
				ProductId:  item.ProductID.String(),
				Quantity:   item.Quantity,
				UnitPrice:  int(item.UnitPrice.Cents()),
				TotalPrice: int(item.UnitPrice.Cents()) * item.Quantity,
			}
		}

		responses[i] = openapi.OrderResponse{
			Id:          order.ID().String(),
			CustomerId:  order.CustomerID().String(),
			Items:       items,
			Status:      openapi.OrderResponseStatus(order.Status()),
			TotalAmount: int(order.Total().Cents()),
			CreatedAt:   order.CreatedAt(),
			UpdatedAt:   order.UpdatedAt(),
		}
	}

	return ctx.JSON(statusCode, responses)
}

// PresentError presents an error response
func (p *OrderPresenter) PresentError(ctx echo.Context, statusCode int, code, message string) error {
	errorResponse := openapi.Error{
		Code:    code,
		Message: message,
	}

	return ctx.JSON(statusCode, errorResponse)
}
