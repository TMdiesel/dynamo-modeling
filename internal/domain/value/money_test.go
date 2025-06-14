package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoney(t *testing.T) {
	t.Run("create money from cents", func(t *testing.T) {
		money, err := NewMoney(150)
		assert.NoError(t, err)
		assert.Equal(t, int64(150), money.Cents())
		assert.Equal(t, 1.50, money.Dollars())
		assert.Equal(t, "$1.50", money.String())
	})

	t.Run("create money from float", func(t *testing.T) {
		money, err := NewMoneyFromFloat(1.50)
		assert.NoError(t, err)
		assert.Equal(t, int64(150), money.Cents())
		assert.Equal(t, 1.50, money.Dollars())
	})

	t.Run("create money from dollars", func(t *testing.T) {
		money, err := NewMoneyFromDollars(5)
		assert.NoError(t, err)
		assert.Equal(t, int64(500), money.Cents())
		assert.Equal(t, 5.00, money.Dollars())
		assert.Equal(t, "$5.00", money.String())
	})

	t.Run("negative money should return error", func(t *testing.T) {
		_, err := NewMoney(-100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")

		_, err = NewMoneyFromFloat(-1.50)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")

		_, err = NewMoneyFromDollars(-5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")
	})

	t.Run("zero money", func(t *testing.T) {
		money, _ := NewMoney(0)
		assert.True(t, money.IsZero())
		assert.False(t, money.IsPositive())
	})

	t.Run("positive money", func(t *testing.T) {
		money, _ := NewMoney(100)
		assert.False(t, money.IsZero())
		assert.True(t, money.IsPositive())
	})

	t.Run("add money", func(t *testing.T) {
		money1, _ := NewMoney(100)
		money2, _ := NewMoney(50)
		result := money1.Add(money2)
		assert.Equal(t, int64(150), result.Cents())
	})

	t.Run("subtract money", func(t *testing.T) {
		money1, _ := NewMoney(150)
		money2, _ := NewMoney(50)
		result, err := money1.Subtract(money2)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), result.Cents())
	})

	t.Run("subtract money resulting in negative should return error", func(t *testing.T) {
		money1, _ := NewMoney(50)
		money2, _ := NewMoney(100)
		_, err := money1.Subtract(money2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "negative amount")
	})

	t.Run("multiply money", func(t *testing.T) {
		money, _ := NewMoney(100)
		result, err := money.Multiply(2.5)
		assert.NoError(t, err)
		assert.Equal(t, int64(250), result.Cents())
	})

	t.Run("multiply money by negative factor should return error", func(t *testing.T) {
		money, _ := NewMoney(100)
		_, err := money.Multiply(-2.0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be negative")
	})

	t.Run("money equals", func(t *testing.T) {
		money1, _ := NewMoney(100)
		money2, _ := NewMoney(100)
		money3, _ := NewMoney(150)

		assert.True(t, money1.Equals(money2))
		assert.False(t, money1.Equals(money3))
	})

	t.Run("money comparison", func(t *testing.T) {
		money1, _ := NewMoney(100)
		money2, _ := NewMoney(150)
		money3, _ := NewMoney(100)

		assert.True(t, money2.GreaterThan(money1))
		assert.False(t, money1.GreaterThan(money2))

		assert.True(t, money1.LessThan(money2))
		assert.False(t, money2.LessThan(money1))

		assert.True(t, money1.GreaterThanOrEqual(money3))
		assert.True(t, money2.GreaterThanOrEqual(money1))
		assert.False(t, money1.GreaterThanOrEqual(money2))

		assert.True(t, money1.LessThanOrEqual(money3))
		assert.True(t, money1.LessThanOrEqual(money2))
		assert.False(t, money2.LessThanOrEqual(money1))
	})

	t.Run("floating point precision", func(t *testing.T) {
		// Test that floating point precision issues are handled correctly
		money, err := NewMoneyFromFloat(1.005) // Should round to 1.00
		assert.NoError(t, err)
		assert.Equal(t, int64(100), money.Cents()) // 100 cents = $1.00
	})
}
