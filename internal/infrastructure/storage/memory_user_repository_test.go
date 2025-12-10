package storage

import (
	"testing"
	"user-management/internal/domain/models"
)

func TestMemoryUserRepository_Save(t *testing.T) {
	repo := NewMemoryUserRepository()

	user := models.NewUser("Test User", "test@example.com", 25, true)

	// Primer guardado
	saved, err := repo.Save(user)
	if err != nil {
		t.Fatalf("Save falló: %v", err)
	}

	if saved.ID == 0 {
		t.Error("Save() debería asignar ID")
	}

	// Verificar que se guardó
	found, err := repo.FindByID(int(saved.ID))
	if err != nil {
		t.Fatalf("FindByID falló: %v", err)
	}

	if found.Name != "Test User" {
		t.Errorf("FindByID() = %v, quiere Test User", found.Name)
	}
}

func TestMemoryUserRepository_FindAll(t *testing.T) {
	repo := NewMemoryUserRepository()

	// Agregar múltiples usuarios
	users := []struct {
		name  string
		email string
		age   int
	}{
		{"User1", "user1@test.com", 20},
		{"User2", "user2@test.com", 30},
		{"User3", "user3@test.com", 40},
	}

	for _, u := range users {
		user := models.NewUser(u.name, u.email, u.age, true)
		_, err := repo.Save(user)
		if err != nil {
			t.Fatalf("Save falló: %v", err)
		}
	}

	allUsers, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll falló: %v", err)
	}
	if len(allUsers) != 3 {
		t.Errorf("FindAll() = %d usuarios, quiere 3", len(allUsers))
	}
}

func TestMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryUserRepository()

	// Test de concurrencia básica
	done := make(chan bool)

	// Goroutine 1
	go func() {
		user := models.NewUser("User1", "user1@test.com", 20, true)
		repo.Save(user)
		done <- true
	}()

	// Goroutine 2
	go func() {
		user := models.NewUser("User2", "user2@test.com", 30, true)
		repo.Save(user)
		done <- true
	}()

	// Esperar ambas goroutines
	<-done
	<-done

	// Verificar que ambos usuarios existen
	allUsers, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll falló: %v", err)
	}

	if len(allUsers) != 2 {
		t.Errorf("Concurrent Save() = %d usuarios, quiere 2", len(allUsers))
	}
}
