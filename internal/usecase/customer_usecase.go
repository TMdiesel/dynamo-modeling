package usecase

import (
	"context"

	"dynamo-modeling/internal/domain"
	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/repository"
	"dynamo-modeling/internal/domain/value"
)

// GetCustomerUseCase handles getting a customer by ID
type GetCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

// GetCustomerCommand represents the input for getting a customer
type GetCustomerCommand struct {
	CustomerID string
}

// NewGetCustomerUseCase creates a new get customer use case
func NewGetCustomerUseCase(customerRepo repository.CustomerRepository) *GetCustomerUseCase {
	return &GetCustomerUseCase{
		customerRepo: customerRepo,
	}
}

// Execute executes the get customer use case
func (uc *GetCustomerUseCase) Execute(ctx context.Context, cmd GetCustomerCommand) (*entity.Customer, error) {
	// 1. 値オブジェクトの作成・バリデーション
	customerID := value.CustomerID(cmd.CustomerID)

	// 2. リポジトリから取得
	customer, err := uc.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, domain.CustomerNotFoundError(cmd.CustomerID)
	}

	return customer, nil
}

// ListCustomersUseCase handles listing customers
type ListCustomersUseCase struct {
	customerRepo repository.CustomerRepository
}

// ListCustomersCommand represents the input for listing customers
type ListCustomersCommand struct {
	Limit int
}

// NewListCustomersUseCase creates a new list customers use case
func NewListCustomersUseCase(customerRepo repository.CustomerRepository) *ListCustomersUseCase {
	return &ListCustomersUseCase{
		customerRepo: customerRepo,
	}
}

// Execute executes the list customers use case
func (uc *ListCustomersUseCase) Execute(ctx context.Context, cmd ListCustomersCommand) ([]*entity.Customer, error) {
	// ビジネスロジック: デフォルトの制限値設定
	if cmd.Limit <= 0 || cmd.Limit > 1000 {
		cmd.Limit = 100 // デフォルト制限
	}

	// リポジトリから取得
	customers, err := uc.customerRepo.ListWithLimit(ctx, &cmd.Limit)
	if err != nil {
		return nil, domain.RepositoryError("failed to list customers", err)
	}

	return customers, nil
}

// UpdateCustomerUseCase handles customer update
type UpdateCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

// UpdateCustomerCommand represents the input for updating a customer
type UpdateCustomerCommand struct {
	CustomerID string
	Name       string
	Email      string
}

// NewUpdateCustomerUseCase creates a new update customer use case
func NewUpdateCustomerUseCase(customerRepo repository.CustomerRepository) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{
		customerRepo: customerRepo,
	}
}

// Execute executes the update customer use case
func (uc *UpdateCustomerUseCase) Execute(ctx context.Context, cmd UpdateCustomerCommand) (*entity.Customer, error) {
	// 1. 値オブジェクトの作成・バリデーション
	customerID := value.CustomerID(cmd.CustomerID)
	email, err := value.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	// 2. 既存の顧客を取得
	customer, err := uc.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, domain.CustomerNotFoundError(cmd.CustomerID)
	}

	// 3. エンティティの更新
	customer.UpdateName(cmd.Name)
	customer.UpdateEmail(email)

	// 4. リポジトリに保存
	err = uc.customerRepo.Save(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

// DeleteCustomerUseCase handles customer deletion
type DeleteCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

// DeleteCustomerCommand represents the input for deleting a customer
type DeleteCustomerCommand struct {
	CustomerID string
}

// NewDeleteCustomerUseCase creates a new delete customer use case
func NewDeleteCustomerUseCase(customerRepo repository.CustomerRepository) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{
		customerRepo: customerRepo,
	}
}

// Execute executes the delete customer use case
func (uc *DeleteCustomerUseCase) Execute(ctx context.Context, cmd DeleteCustomerCommand) error {
	// 1. 値オブジェクトの作成・バリデーション
	customerID := value.CustomerID(cmd.CustomerID)

	// 2. 既存の顧客を確認
	customer, err := uc.customerRepo.FindByID(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return domain.CustomerNotFoundError(cmd.CustomerID)
	}

	// 3. ビジネスルール: 注文がある顧客は削除できない（実装予定）
	// TODO: OrderRepository で該当顧客の注文をチェック

	// 4. リポジトリから削除
	return uc.customerRepo.Delete(ctx, customerID)
}
