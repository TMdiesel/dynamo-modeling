package usecase

import (
	"context"

	"dynamo-modeling/internal/domain"
	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/repository"
	"dynamo-modeling/internal/domain/value"
)

// CreateProductUseCase handles product creation business logic
type CreateProductUseCase struct {
	productRepo repository.ProductRepository
}

// CreateProductCommand represents the input for creating a product
type CreateProductCommand struct {
	Name        string
	Description string
	Price       int64
	Stock       int
}

// NewCreateProductUseCase creates a new create product use case
func NewCreateProductUseCase(productRepo repository.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{
		productRepo: productRepo,
	}
}

// Execute executes the create product use case
func (uc *CreateProductUseCase) Execute(ctx context.Context, cmd CreateProductCommand) (*entity.Product, error) {
	// 1. 値オブジェクトの作成・バリデーション
	price, err := value.NewMoney(cmd.Price)
	if err != nil {
		return nil, domain.InvalidInputError("invalid price format")
	}

	// 2. 新しいProduct IDを生成
	productID := value.GenerateProductID()

	// 3. エンティティ作成
	product, err := entity.NewProduct(productID, cmd.Name, cmd.Description, price, cmd.Stock)
	if err != nil {
		return nil, domain.InvalidInputError("failed to create product: " + err.Error())
	}

	// 4. リポジトリに保存
	err = uc.productRepo.Save(ctx, product)
	if err != nil {
		return nil, domain.RepositoryError("failed to save product", err)
	}

	return product, nil
}

// GetProductUseCase handles getting a product by ID
type GetProductUseCase struct {
	productRepo repository.ProductRepository
}

// GetProductCommand represents the input for getting a product
type GetProductCommand struct {
	ProductID string
}

// NewGetProductUseCase creates a new get product use case
func NewGetProductUseCase(productRepo repository.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{
		productRepo: productRepo,
	}
}

// Execute executes the get product use case
func (uc *GetProductUseCase) Execute(ctx context.Context, cmd GetProductCommand) (*entity.Product, error) {
	// 1. 値オブジェクトの作成・バリデーション
	productID := value.ProductID(cmd.ProductID)

	// 2. リポジトリから取得
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, domain.NewDomainError("PRODUCT_NOT_FOUND", "Product not found", nil)
	}

	return product, nil
}

// ListProductsUseCase handles listing products
type ListProductsUseCase struct {
	productRepo repository.ProductRepository
}

// ListProductsCommand represents the input for listing products
type ListProductsCommand struct {
	Limit int
}

// NewListProductsUseCase creates a new list products use case
func NewListProductsUseCase(productRepo repository.ProductRepository) *ListProductsUseCase {
	return &ListProductsUseCase{
		productRepo: productRepo,
	}
}

// Execute executes the list products use case
func (uc *ListProductsUseCase) Execute(ctx context.Context, cmd ListProductsCommand) ([]*entity.Product, error) {
	// ビジネスロジック: デフォルトの制限値設定
	if cmd.Limit <= 0 || cmd.Limit > 1000 {
		cmd.Limit = 100 // デフォルト制限
	}

	// リポジトリから取得
	products, _, err := uc.productRepo.FindAll(ctx, cmd.Limit, nil)
	if err != nil {
		return nil, domain.RepositoryError("failed to list products", err)
	}

	return products, nil
}

// UpdateProductUseCase handles product update
type UpdateProductUseCase struct {
	productRepo repository.ProductRepository
}

// UpdateProductCommand represents the input for updating a product
type UpdateProductCommand struct {
	ProductID   string
	Name        string
	Description string
	Price       int64
	Stock       int
}

// NewUpdateProductUseCase creates a new update product use case
func NewUpdateProductUseCase(productRepo repository.ProductRepository) *UpdateProductUseCase {
	return &UpdateProductUseCase{
		productRepo: productRepo,
	}
}

// Execute executes the update product use case
func (uc *UpdateProductUseCase) Execute(ctx context.Context, cmd UpdateProductCommand) (*entity.Product, error) {
	// 1. 値オブジェクトの作成・バリデーション
	productID := value.ProductID(cmd.ProductID)
	price, err := value.NewMoney(cmd.Price)
	if err != nil {
		return nil, domain.InvalidInputError("invalid price format")
	}

	// 2. 既存の商品を取得
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, domain.NewDomainError("PRODUCT_NOT_FOUND", "Product not found", nil)
	}

	// 3. エンティティの更新
	product.UpdateName(cmd.Name)
	product.UpdateDescription(cmd.Description)
	product.UpdatePrice(price)
	err = product.UpdateStock(cmd.Stock)
	if err != nil {
		return nil, domain.InvalidInputError("invalid stock value: " + err.Error())
	}

	// 4. リポジトリに保存
	err = uc.productRepo.Save(ctx, product)
	if err != nil {
		return nil, domain.RepositoryError("failed to update product", err)
	}

	return product, nil
}

// DeleteProductUseCase handles product deletion
type DeleteProductUseCase struct {
	productRepo repository.ProductRepository
}

// DeleteProductCommand represents the input for deleting a product
type DeleteProductCommand struct {
	ProductID string
}

// NewDeleteProductUseCase creates a new delete product use case
func NewDeleteProductUseCase(productRepo repository.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{
		productRepo: productRepo,
	}
}

// Execute executes the delete product use case
func (uc *DeleteProductUseCase) Execute(ctx context.Context, cmd DeleteProductCommand) error {
	// 1. 値オブジェクトの作成・バリデーション
	productID := value.ProductID(cmd.ProductID)

	// 2. 商品の存在確認
	exists, err := uc.productRepo.Exists(ctx, productID)
	if err != nil {
		return domain.RepositoryError("failed to check product existence", err)
	}

	if !exists {
		return domain.NewDomainError("PRODUCT_NOT_FOUND", "Product not found", nil)
	}

	// 3. ビジネスルール: 注文に含まれている商品は削除できない（将来実装）
	// TODO: OrderRepository で該当商品を参照している注文をチェック

	// 4. リポジトリから削除
	return uc.productRepo.Delete(ctx, productID)
}
