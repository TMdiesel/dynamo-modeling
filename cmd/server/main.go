package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"dynamo-modeling/internal/adapter/controller"
	"dynamo-modeling/internal/adapter/openapi"
	"dynamo-modeling/internal/adapter/presenter"
	"dynamo-modeling/internal/adapter/repository"
	"dynamo-modeling/internal/handler"
	"dynamo-modeling/internal/infrastructure"
	"dynamo-modeling/internal/usecase"
)

func main() {
	slog.Info("DynamoDB + Clean Architecture Online Shop API")
	slog.Info("Starting server...")

	// DynamoDB設定
	dbConfig := infrastructure.DynamoDBConfig{
		Region:    "ap-northeast-1",
		Endpoint:  "http://localhost:8000", // DynamoDB Local
		TableName: "OnlineShop",
	}

	// DynamoDBクライアント初期化
	ctx := context.Background()
	dbClient, err := infrastructure.NewDynamoDBClient(ctx, dbConfig)
	if err != nil {
		slog.Error("Failed to initialize DynamoDB client", "error", err)
		os.Exit(1)
	}

	// Repository層を初期化
	customerRepo := repository.NewDynamoCustomerRepository(dbClient)
	productRepo := repository.NewDynamoProductRepository(dbClient)
	orderRepo := repository.NewDynamoOrderRepository(dbClient)

	// UseCase層を初期化
	// Customer UseCases
	createCustomerUseCase := usecase.NewCreateCustomerUseCase(customerRepo)
	getCustomerUseCase := usecase.NewGetCustomerUseCase(customerRepo)
	listCustomersUseCase := usecase.NewListCustomersUseCase(customerRepo)
	updateCustomerUseCase := usecase.NewUpdateCustomerUseCase(customerRepo)
	deleteCustomerUseCase := usecase.NewDeleteCustomerUseCase(customerRepo)

	// Product UseCases
	createProductUseCase := usecase.NewCreateProductUseCase(productRepo)
	getProductUseCase := usecase.NewGetProductUseCase(productRepo)
	listProductsUseCase := usecase.NewListProductsUseCase(productRepo)
	updateProductUseCase := usecase.NewUpdateProductUseCase(productRepo)
	deleteProductUseCase := usecase.NewDeleteProductUseCase(productRepo)

	// Order UseCases
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepo, customerRepo, productRepo)
	getOrderUseCase := usecase.NewGetOrderUseCase(orderRepo)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepo)
	updateOrderStatusUseCase := usecase.NewUpdateOrderStatusUseCase(orderRepo)

	// Presenter層を初期化
	customerPresenter := presenter.NewCustomerPresenter()
	productPresenter := presenter.NewProductPresenter()
	orderPresenter := presenter.NewOrderPresenter()

	// Controller層を初期化
	customerController := controller.NewCustomerController(
		createCustomerUseCase,
		getCustomerUseCase,
		listCustomersUseCase,
		updateCustomerUseCase,
		deleteCustomerUseCase,
		customerPresenter,
	)

	productController := controller.NewProductController(
		createProductUseCase,
		getProductUseCase,
		listProductsUseCase,
		updateProductUseCase,
		deleteProductUseCase,
		productPresenter,
	)

	orderController := controller.NewOrderController(
		createOrderUseCase,
		getOrderUseCase,
		listOrdersUseCase,
		updateOrderStatusUseCase,
		orderPresenter,
	)

	// Handler層を初期化
	apiHandler := handler.NewAPIHandler(customerController, productController, orderController)

	// Echoサーバー作成
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// OpenAPIハンドラーを登録
	openapi.RegisterHandlers(e, apiHandler)

	// サーバー設定とグレースフルシャットダウン

	// グレースフルシャットダウン設定
	go func() {
		slog.Info("Server starting on :8080")
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server shutting down...")

	// グレースフルシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited")
}
