package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		mockValid bool
		wantErr   bool
	}{
		{"valid", "test@example.com", true, false},
		{"invalid", "invalid", false, true},
		{"empty", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reset := setEmailValidator(func(email string) bool {
				return tt.mockValid
			})
			defer reset()

			emailVO, err := NewEmail(tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidEmail, err)
				assert.Empty(t, emailVO.Value())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, emailVO.Value())
			}
		})
	}
}
