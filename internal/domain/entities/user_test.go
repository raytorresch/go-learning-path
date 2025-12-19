package entities_test

import (
	"testing"
	"user-management/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("success - valid user", func(t *testing.T) {
		user, err := entities.NewUser("test@example.com", "John Doe", 0, "SecurePass123!")

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John Doe", user.Name)
		assert.False(t, user.ID == uuid.Nil)
	})

	t.Run("failure - invalid email", func(t *testing.T) {
		user, err := entities.NewUser("invalid-email", "John Doe", 0, "SecurePass123!")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, entities.ErrInvalidEmail)
	})

	// Table-driven tests (patrón industrial)
	tests := []struct {
		name     string
		email    string
		password string
		age      int
		expected error
	}{
		{"empty email", "", "pass123", 0, entities.ErrInvalidEmail},
		{"weak password", "test@mail.com", "123", 0, entities.ErrWeakPassword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entities.NewUser(tt.email, "Test", tt.age, tt.password)
			assert.ErrorIs(t, err, tt.expected)
		})
	}
}
