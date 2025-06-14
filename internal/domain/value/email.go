package value

import (
	"fmt"
	"regexp"
	"strings"
)

// Email represents a valid email address
type Email struct {
	value string
}

// emailRegex is a simple email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email with validation
func NewEmail(email string) (Email, error) {
	trimmed := strings.TrimSpace(email)
	if trimmed == "" {
		return Email{}, fmt.Errorf("email cannot be empty")
	}

	if !emailRegex.MatchString(trimmed) {
		return Email{}, fmt.Errorf("invalid email format: %s", trimmed)
	}

	return Email{value: strings.ToLower(trimmed)}, nil
}

// String returns the string representation of Email
func (e Email) String() string {
	return e.value
}

// Value returns the raw email value
func (e Email) Value() string {
	return e.value
}

// IsEmpty checks if the Email is empty
func (e Email) IsEmpty() bool {
	return e.value == ""
}

// Equals compares two Email values
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Domain returns the domain part of the email
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart returns the local part of the email (before @)
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}
