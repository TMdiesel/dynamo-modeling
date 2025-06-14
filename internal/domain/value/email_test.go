package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	t.Run("valid email", func(t *testing.T) {
		email, err := NewEmail("test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", email.String())
		assert.Equal(t, "test@example.com", email.Value())
		assert.False(t, email.IsEmpty())
	})

	t.Run("email should be lowercased", func(t *testing.T) {
		email, err := NewEmail("Test@Example.Com")
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", email.String())
	})

	t.Run("empty email should return error", func(t *testing.T) {
		_, err := NewEmail("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("whitespace only email should return error", func(t *testing.T) {
		_, err := NewEmail("   ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("invalid email format should return error", func(t *testing.T) {
		testCases := []string{
			"invalid-email",
			"@example.com",
			"test@",
			"test@example",
		}

		for _, tc := range testCases {
			_, err := NewEmail(tc)
			assert.Error(t, err, "expected error for email: %s", tc)
			if err != nil {
				assert.Contains(t, err.Error(), "invalid email format")
			}
		}
	})

	t.Run("email equals", func(t *testing.T) {
		email1, _ := NewEmail("test@example.com")
		email2, _ := NewEmail("test@example.com")
		email3, _ := NewEmail("different@example.com")

		assert.True(t, email1.Equals(email2))
		assert.False(t, email1.Equals(email3))
	})

	t.Run("email domain", func(t *testing.T) {
		email, _ := NewEmail("test@example.com")
		assert.Equal(t, "example.com", email.Domain())
	})

	t.Run("email local part", func(t *testing.T) {
		email, _ := NewEmail("test@example.com")
		assert.Equal(t, "test", email.LocalPart())
	})

	t.Run("empty email check", func(t *testing.T) {
		var email Email
		assert.True(t, email.IsEmpty())
	})

	t.Run("email with whitespace should be trimmed", func(t *testing.T) {
		email, err := NewEmail("  test@example.com  ")
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", email.String())
	})
}
