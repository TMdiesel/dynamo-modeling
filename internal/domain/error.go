package domain

import "fmt"

// DomainError represents domain-specific errors
type DomainError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the error unwrapping interface
func (e *DomainError) Unwrap() error {
	return e.Err
}

// Common domain error codes
const (
	ErrCodeCustomerNotFound      = "CUSTOMER_NOT_FOUND"
	ErrCodeCustomerAlreadyExists = "CUSTOMER_ALREADY_EXISTS"
	ErrCodeInvalidInput          = "INVALID_INPUT"
	ErrCodeRepositoryError       = "REPOSITORY_ERROR"
)

// NewDomainError creates a new domain error
func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// CustomerNotFoundError creates a customer not found error
func CustomerNotFoundError(customerID string) *DomainError {
	return NewDomainError(
		ErrCodeCustomerNotFound,
		fmt.Sprintf("Customer with ID %s not found", customerID),
		nil,
	)
}

// CustomerAlreadyExistsError creates a customer already exists error
func CustomerAlreadyExistsError(email string) *DomainError {
	return NewDomainError(
		ErrCodeCustomerAlreadyExists,
		fmt.Sprintf("Customer with email %s already exists", email),
		nil,
	)
}

// InvalidInputError creates an invalid input error
func InvalidInputError(message string) *DomainError {
	return NewDomainError(
		ErrCodeInvalidInput,
		message,
		nil,
	)
}

// RepositoryError creates a repository error
func RepositoryError(message string, err error) *DomainError {
	return NewDomainError(
		ErrCodeRepositoryError,
		message,
		err,
	)
}
