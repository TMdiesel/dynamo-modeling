package value

import (
	"fmt"
	"math"
)

// Money represents a monetary amount in cents to avoid floating point precision issues
type Money struct {
	cents int64
}

// NewMoney creates a new Money value from cents
func NewMoney(cents int64) (Money, error) {
	if cents < 0 {
		return Money{}, fmt.Errorf("money amount cannot be negative: %d", cents)
	}
	return Money{cents: cents}, nil
}

// NewMoneyFromFloat creates a new Money value from a float amount (in dollars)
func NewMoneyFromFloat(amount float64) (Money, error) {
	if amount < 0 {
		return Money{}, fmt.Errorf("money amount cannot be negative: %.2f", amount)
	}

	// Convert to cents and round to avoid floating point precision issues
	cents := int64(math.Round(amount * 100))
	return Money{cents: cents}, nil
}

// NewMoneyFromDollars creates a new Money value from dollar amount
func NewMoneyFromDollars(dollars int64) (Money, error) {
	if dollars < 0 {
		return Money{}, fmt.Errorf("money amount cannot be negative: %d", dollars)
	}
	return Money{cents: dollars * 100}, nil
}

// Cents returns the amount in cents
func (m Money) Cents() int64 {
	return m.cents
}

// Dollars returns the amount in dollars as a float
func (m Money) Dollars() float64 {
	return float64(m.cents) / 100.0
}

// String returns a formatted string representation of the money
func (m Money) String() string {
	return fmt.Sprintf("$%.2f", m.Dollars())
}

// IsZero checks if the money amount is zero
func (m Money) IsZero() bool {
	return m.cents == 0
}

// IsPositive checks if the money amount is positive
func (m Money) IsPositive() bool {
	return m.cents > 0
}

// Add adds two Money values
func (m Money) Add(other Money) Money {
	return Money{cents: m.cents + other.cents}
}

// Subtract subtracts another Money value from this one
func (m Money) Subtract(other Money) (Money, error) {
	result := m.cents - other.cents
	if result < 0 {
		return Money{}, fmt.Errorf("subtraction would result in negative amount")
	}
	return Money{cents: result}, nil
}

// Multiply multiplies the money by a factor
func (m Money) Multiply(factor float64) (Money, error) {
	if factor < 0 {
		return Money{}, fmt.Errorf("multiplication factor cannot be negative: %.2f", factor)
	}

	result := int64(math.Round(float64(m.cents) * factor))
	return Money{cents: result}, nil
}

// Equals compares two Money values
func (m Money) Equals(other Money) bool {
	return m.cents == other.cents
}

// GreaterThan checks if this money is greater than another
func (m Money) GreaterThan(other Money) bool {
	return m.cents > other.cents
}

// LessThan checks if this money is less than another
func (m Money) LessThan(other Money) bool {
	return m.cents < other.cents
}

// GreaterThanOrEqual checks if this money is greater than or equal to another
func (m Money) GreaterThanOrEqual(other Money) bool {
	return m.cents >= other.cents
}

// LessThanOrEqual checks if this money is less than or equal to another
func (m Money) LessThanOrEqual(other Money) bool {
	return m.cents <= other.cents
}
