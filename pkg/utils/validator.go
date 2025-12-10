package utils

import "strings"

// ValidateEmail es exportada y reusable
func ValidateEmail(email string) bool {
	// Lógica de validación compleja
	return len(email) > 3 && strings.Contains(email, "@")
}

// ValidateAge es exportada
func ValidateAge(age int) bool {
	return age >= 0 && age <= 120
}
