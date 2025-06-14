package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/guregu/dynamo/v2"

	"dynamo-modeling/internal/domain/entity"
	"dynamo-modeling/internal/domain/value"
	"dynamo-modeling/internal/infrastructure"
)

// DynamoCustomerRepository implements CustomerRepository using DynamoDB
type DynamoCustomerRepository struct {
	client *infrastructure.DynamoDBClient
}

// NewDynamoCustomerRepository creates a new DynamoDB customer repository
func NewDynamoCustomerRepository(client *infrastructure.DynamoDBClient) *DynamoCustomerRepository {
	return &DynamoCustomerRepository{
		client: client,
	}
}

// CustomerItem represents a customer item in DynamoDB
type CustomerItem struct {
	PK        string    `dynamo:"PK"`        // CUSTOMER#{CustomerID}
	SK        string    `dynamo:"SK"`        // CUSTOMER#{CustomerID}
	GSI1PK    string    `dynamo:"GSI1PK"`    // EMAIL#{Email}
	GSI1SK    string    `dynamo:"GSI1SK"`    // CUSTOMER#{CustomerID}
	Type      string    `dynamo:"Type"`      // "CUSTOMER"
	ID        string    `dynamo:"ID"`        // CustomerID
	Email     string    `dynamo:"Email"`     // Email address
	Name      string    `dynamo:"Name"`      // Customer name
	CreatedAt time.Time `dynamo:"CreatedAt"` // Creation timestamp
	UpdatedAt time.Time `dynamo:"UpdatedAt"` // Last update timestamp
}

// ToEntity converts CustomerItem to Customer entity
func (item *CustomerItem) ToEntity() (*entity.Customer, error) {
	customerID, err := value.NewCustomerID(item.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID: %w", err)
	}

	email, err := value.NewEmail(item.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	customer := entity.NewCustomer(customerID, email, item.Name)

	// Set timestamps using reflection or recreate entity with timestamps
	// For now, we'll create a new instance and trust the timestamps are correct
	return customer, nil
}

// FromEntity converts Customer entity to CustomerItem
func CustomerItemFromEntity(customer *entity.Customer) *CustomerItem {
	customerID := customer.ID().String()

	return &CustomerItem{
		PK:        fmt.Sprintf("CUSTOMER#%s", customerID),
		SK:        fmt.Sprintf("CUSTOMER#%s", customerID),
		GSI1PK:    fmt.Sprintf("EMAIL#%s", customer.Email().String()),
		GSI1SK:    fmt.Sprintf("CUSTOMER#%s", customerID),
		Type:      "CUSTOMER",
		ID:        customerID,
		Email:     customer.Email().String(),
		Name:      customer.Name(),
		CreatedAt: customer.CreatedAt(),
		UpdatedAt: customer.UpdatedAt(),
	}
}

// Save creates or updates a customer
func (r *DynamoCustomerRepository) Save(ctx context.Context, customer *entity.Customer) error {
	slog.Info("Saving customer", "customerID", customer.ID().String())

	// Check if email is already taken by another customer
	if err := r.checkEmailUniqueness(ctx, customer); err != nil {
		return err
	}

	item := CustomerItemFromEntity(customer)
	table := r.client.GetTable()

	err := table.Put(item).Run(ctx)
	if err != nil {
		slog.Error("Failed to save customer", "customerID", customer.ID().String(), "error", err)
		return fmt.Errorf("failed to save customer: %w", err)
	}

	slog.Info("Customer saved successfully", "customerID", customer.ID().String())
	return nil
}

// FindByID retrieves a customer by their ID
func (r *DynamoCustomerRepository) FindByID(ctx context.Context, id value.CustomerID) (*entity.Customer, error) {
	slog.Info("Finding customer by ID", "customerID", id.String())

	pk := fmt.Sprintf("CUSTOMER#%s", id.String())
	sk := fmt.Sprintf("CUSTOMER#%s", id.String())

	var item CustomerItem
	table := r.client.GetTable()

	err := table.Get("PK", pk).Range("SK", dynamo.Equal, sk).One(ctx, &item)
	if err != nil {
		if err == dynamo.ErrNotFound {
			slog.Warn("Customer not found", "customerID", id.String())
			return nil, fmt.Errorf("customer not found")
		}
		slog.Error("Failed to find customer", "customerID", id.String(), "error", err)
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	customer, err := item.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed to convert item to entity: %w", err)
	}

	slog.Info("Customer found successfully", "customerID", id.String())
	return customer, nil
}

// FindByEmail retrieves a customer by their email address using GSI1
func (r *DynamoCustomerRepository) FindByEmail(ctx context.Context, email value.Email) (*entity.Customer, error) {
	slog.Info("Finding customer by email", "email", email.String())

	gsi1pk := fmt.Sprintf("EMAIL#%s", email.String())

	var item CustomerItem
	table := r.client.GetTable()

	err := table.Get("GSI1PK", gsi1pk).Index("GSI1").One(ctx, &item)
	if err != nil {
		if err == dynamo.ErrNotFound {
			slog.Warn("Customer not found by email", "email", email.String())
			return nil, fmt.Errorf("customer not found")
		}
		slog.Error("Failed to find customer by email", "email", email.String(), "error", err)
		return nil, fmt.Errorf("failed to find customer by email: %w", err)
	}

	customer, err := item.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("failed to convert item to entity: %w", err)
	}

	slog.Info("Customer found by email successfully", "email", email.String())
	return customer, nil
}

// Delete removes a customer by their ID
func (r *DynamoCustomerRepository) Delete(ctx context.Context, id value.CustomerID) error {
	slog.Info("Deleting customer", "customerID", id.String())

	pk := fmt.Sprintf("CUSTOMER#%s", id.String())
	sk := fmt.Sprintf("CUSTOMER#%s", id.String())

	table := r.client.GetTable()

	err := table.Delete("PK", pk).Range("SK", sk).Run(ctx)
	if err != nil {
		slog.Error("Failed to delete customer", "customerID", id.String(), "error", err)
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	slog.Info("Customer deleted successfully", "customerID", id.String())
	return nil
}

// Exists checks if a customer exists by their ID
func (r *DynamoCustomerRepository) Exists(ctx context.Context, id value.CustomerID) (bool, error) {
	slog.Info("Checking if customer exists", "customerID", id.String())

	_, err := r.FindByID(ctx, id)
	if err != nil {
		if err.Error() == "customer not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// checkEmailUniqueness ensures that the email is not already taken by another customer
func (r *DynamoCustomerRepository) checkEmailUniqueness(ctx context.Context, customer *entity.Customer) error {
	existingCustomer, err := r.FindByEmail(ctx, customer.Email())
	if err != nil {
		if err.Error() == "customer not found" {
			// Email is available
			return nil
		}
		// Other error occurred
		return err
	}

	// Check if the existing customer is the same as the one being saved
	if existingCustomer.ID() != customer.ID() {
		return fmt.Errorf("email already taken by another customer")
	}

	return nil
}
