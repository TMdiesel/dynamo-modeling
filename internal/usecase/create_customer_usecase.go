package usecase

import (
	"context"

	"dynamo-modeling/internal/domain"
	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/repository"
	"dynamo-modeling/internal/domain/value"
)

// CreateCustomerUseCase handles customer creation business logic
type CreateCustomerUseCase struct {
	customerRepo repository.CustomerRepository
}

// CreateCustomerCommand represents the input for creating a customer
type CreateCustomerCommand struct {
	Name  string
	Email string
}

// NewCreateCustomerUseCase creates a new create customer use case
func NewCreateCustomerUseCase(customerRepo repository.CustomerRepository) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{
		customerRepo: customerRepo,
	}
}

// Execute executes the create customer use case
func (uc *CreateCustomerUseCase) Execute(ctx context.Context, cmd CreateCustomerCommand) (*entity.Customer, error) {
	// 1. 値オブジェクトの作成・バリデーション
	email, err := value.NewEmail(cmd.Email)
	if err != nil {
		return nil, domain.InvalidInputError("invalid email format")
	}

	// 2. 新しいCustomer IDを生成
	customerID := value.GenerateCustomerID()

	// 3. エンティティ作成
	customer := entity.NewCustomer(customerID, email, cmd.Name)

	// 4. ビジネスルール: 同じメールアドレスの顧客が存在しないかチェック
	existingCustomer, err := uc.customerRepo.FindByEmail(ctx, email)
	if err == nil && existingCustomer != nil {
		return nil, domain.CustomerAlreadyExistsError(cmd.Email)
	}

	// 5. リポジトリに保存
	err = uc.customerRepo.Save(ctx, customer)
	if err != nil {
		return nil, domain.RepositoryError("failed to save customer", err)
	}

	return customer, nil
}
