package usecase_test

import (
	"context"
	"testing"

	"dynamo-modeling/internal/domain"
	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
	"dynamo-modeling/internal/usecase"
)

// MockCustomerRepository implements CustomerRepository for testing
type MockCustomerRepository struct {
	customers  map[string]*entity.Customer
	emailIndex map[string]*entity.Customer
}

func NewMockCustomerRepository() *MockCustomerRepository {
	return &MockCustomerRepository{
		customers:  make(map[string]*entity.Customer),
		emailIndex: make(map[string]*entity.Customer),
	}
}

func (m *MockCustomerRepository) Save(ctx context.Context, customer *entity.Customer) error {
	m.customers[customer.ID().String()] = customer
	m.emailIndex[customer.Email().String()] = customer
	return nil
}

func (m *MockCustomerRepository) FindByID(ctx context.Context, id value.CustomerID) (*entity.Customer, error) {
	customer, exists := m.customers[id.String()]
	if !exists {
		return nil, domain.CustomerNotFoundError(id.String())
	}
	return customer, nil
}

func (m *MockCustomerRepository) FindByEmail(ctx context.Context, email value.Email) (*entity.Customer, error) {
	customer, exists := m.emailIndex[email.String()]
	if !exists {
		return nil, domain.CustomerNotFoundError(email.String())
	}
	return customer, nil
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id value.CustomerID) error {
	customer, exists := m.customers[id.String()]
	if !exists {
		return domain.CustomerNotFoundError(id.String())
	}
	delete(m.customers, id.String())
	delete(m.emailIndex, customer.Email().String())
	return nil
}

func (m *MockCustomerRepository) Exists(ctx context.Context, id value.CustomerID) (bool, error) {
	_, exists := m.customers[id.String()]
	return exists, nil
}

func (m *MockCustomerRepository) ListWithLimit(ctx context.Context, limit *int) ([]*entity.Customer, error) {
	customers := make([]*entity.Customer, 0, len(m.customers))
	count := 0
	for _, customer := range m.customers {
		if limit != nil && count >= *limit {
			break
		}
		customers = append(customers, customer)
		count++
	}
	return customers, nil
}

func TestCreateCustomerUseCase_Success(t *testing.T) {
	// Arrange
	repo := NewMockCustomerRepository()
	uc := usecase.NewCreateCustomerUseCase(repo)
	ctx := context.Background()

	cmd := usecase.CreateCustomerCommand{
		Name:  "Test Customer",
		Email: "test@example.com",
	}

	// Act
	customer, err := uc.Execute(ctx, cmd)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if customer == nil {
		t.Fatal("Expected customer to be created")
	}

	if customer.Name() != cmd.Name {
		t.Errorf("Expected name %s, got %s", cmd.Name, customer.Name())
	}

	if customer.Email().String() != cmd.Email {
		t.Errorf("Expected email %s, got %s", cmd.Email, customer.Email().String())
	}
}

func TestCreateCustomerUseCase_DuplicateEmail(t *testing.T) {
	// Arrange
	repo := NewMockCustomerRepository()
	uc := usecase.NewCreateCustomerUseCase(repo)
	ctx := context.Background()

	// Create first customer
	email, _ := value.NewEmail("test@example.com")
	existingCustomer := entity.NewCustomer(value.GenerateCustomerID(), email, "Existing Customer")
	repo.Save(ctx, existingCustomer)

	cmd := usecase.CreateCustomerCommand{
		Name:  "Test Customer",
		Email: "test@example.com", // Same email
	}

	// Act
	customer, err := uc.Execute(ctx, cmd)

	// Assert
	if err == nil {
		t.Fatal("Expected error for duplicate email")
	}

	if customer != nil {
		t.Fatal("Expected no customer to be created")
	}

	domainErr, ok := err.(*domain.DomainError)
	if !ok {
		t.Fatalf("Expected DomainError, got %T", err)
	}

	if domainErr.Code != domain.ErrCodeCustomerAlreadyExists {
		t.Errorf("Expected error code %s, got %s", domain.ErrCodeCustomerAlreadyExists, domainErr.Code)
	}
}

func TestCreateCustomerUseCase_InvalidEmail(t *testing.T) {
	// Arrange
	repo := NewMockCustomerRepository()
	uc := usecase.NewCreateCustomerUseCase(repo)
	ctx := context.Background()

	cmd := usecase.CreateCustomerCommand{
		Name:  "Test Customer",
		Email: "invalid-email", // Invalid email format
	}

	// Act
	customer, err := uc.Execute(ctx, cmd)

	// Assert
	if err == nil {
		t.Fatal("Expected error for invalid email")
	}

	if customer != nil {
		t.Fatal("Expected no customer to be created")
	}
}
