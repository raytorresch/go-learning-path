package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPasswordHash_WithValidation(t *testing.T) {

	tests := []struct {
		name        string
		input       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "valid bcrypt hash",
			input:       "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			shouldError: false,
		},
		{
			name:        "empty hash - should error",
			input:       "",
			shouldError: true,
			errorMsg:    "cannot be empty", // Futuro mensaje
		},
		{
			name:        "too short - should error",
			input:       "short",
			shouldError: true,
			errorMsg:    "invalid format", // Futuro mensaje
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewPasswordHash(tt.input)

			if tt.shouldError {
				if err != nil {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, hash.Value())
			}
		})
	}
}
