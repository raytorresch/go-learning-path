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
	isValidEmail    = defaultIsValidEmail
)

func defaultIsValidEmail(email string) bool {
	// Importamos utils solo aquí
	return utils.ValidateEmail(email)
}

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

func setEmailValidator(fn func(string) bool) func() {
	original := isValidEmail
	isValidEmail = fn
	return func() { isValidEmail = original }
}
