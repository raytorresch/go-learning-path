package models

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	// Table-driven test
	tests := []struct {
		name    string
		email   string
		age     int
		active  bool
		wantErr bool
	}{
		{"Juan", "juan@test.com", 30, true, false},
		{"", "test@test.com", 25, false, false},     // ¿Nombre vacío permitido?
		{"María", "invalid-email", 20, true, true},  // Email inválido
		{"Pedro", "pedro@test.com", -5, true, true}, // Edad negativa
		{"Ana", "ana@test.com", 150, true, true},    // Edad excesiva
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.name, tt.email, tt.age, tt.active)

			if tt.wantErr {
				// Verificar que email inválido fue rechazado
				if err := user.SetEmail(tt.email); err == nil && tt.email == "invalid-email" {
					t.Errorf("NewUser(%q, %q, %d) debería fallar",
						tt.name, tt.email, tt.age)
				}
			} else {
				// Verificaciones positivas
				if user.Name != tt.name {
					t.Errorf("Nombre = %q, quiere %q", user.Name, tt.name)
				}
				if user.Age != tt.age {
					t.Errorf("Edad = %d, quiere %d", user.Age, tt.age)
				}
			}
		})
	}
}

func TestUser_EmailGetterSetter(t *testing.T) {
	user := NewUser("Test", "old@test.com", 25, true)

	// Test getter
	if user.Email() != "old@test.com" {
		t.Errorf("Email() = %q, quiere %q", user.Email(), "old@test.com")
	}

	// Test setter válido
	if err := user.SetEmail("new@test.com"); err != nil {
		t.Errorf("SetEmail válido falló: %v", err)
	}

	// Test setter inválido
	if err := user.SetEmail("invalid"); err == nil {
		t.Errorf("SetEmail inválido debería fallar")
	}
}
