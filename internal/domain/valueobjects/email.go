package valueobjects

import (
	"errors"
	"user-management/pkg/utils"
)

type Email struct {
	value string
}

var (
	ErrInvalidEmail = errors.New("invalid email format")
)

func NewEmail(value string) (Email, error) {
	// Validación de negocio
	if !isValidEmail(value) {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: value}, nil
}

func (e Email) Value() string {
	return e.value
}

func isValidEmail(email string) bool {
	// Implementar lógica de validación de email
	return utils.ValidateEmail(email)
}
