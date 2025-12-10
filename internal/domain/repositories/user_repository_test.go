package repositories

import (
	"fmt"
	"testing"
	"user-management/internal/domain/models"
)

// MockUserRepository para testing
type MockUserRepository struct {
	users       map[int]*models.User
	saveError   error
	deleteError error
}

// Delete implements UserRepository.
func (m *MockUserRepository) Delete(id int) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	delete(m.users, id)
	return nil
}

// FindAll implements UserRepository.
func (m *MockUserRepository) FindAll() ([]*models.User, error) {
	users := make([]*models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

// Update implements UserRepository.
func (m *MockUserRepository) Update(user *models.User) (*models.User, error) {
	if _, exists := m.users[int(user.ID)]; !exists {
		return nil, fmt.Errorf("not found")
	}
	m.users[int(user.ID)] = user
	return user, nil
}

func (m *MockUserRepository) Save(user *models.User) (*models.User, error) {
	if m.saveError != nil {
		return nil, m.saveError
	}
	m.users[int(user.ID)] = user
	return user, nil
}

func (m *MockUserRepository) FindByID(id int) (*models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, fmt.Errorf("not found")
	}
	return user, nil
}

// Implementa otros métodos...

func TestUserRepositoryInterface(t *testing.T) {
	// Verificar que mock implementa la interfaz
	var _ UserRepository = (*MockUserRepository)(nil)

	mockRepo := &MockUserRepository{
		users: make(map[int]*models.User),
	}

	user := models.NewUser("Test", "test@test.com", 30, true)

	// Test Save
	saved, err := mockRepo.Save(user)
	if err != nil {
		t.Fatalf("Save falló: %v", err)
	}

	if saved.Name != "Test" {
		t.Errorf("Save() = %v, quiere Test", saved.Name)
	}

	// Test FindByID
	found, err := mockRepo.FindByID(int(user.ID))
	if err != nil {
		t.Fatalf("FindByID falló: %v", err)
	}

	if found != saved {
		t.Errorf("FindByID() = %v, quiere %v", found, saved)
	}

	//found all
	all, err := mockRepo.FindAll()
	if err != nil {
		t.Fatalf("FindAll falló: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("FindAll() = %d, quiere 1", len(all))
	}

	// Test Update
	user.Name = "Updated"
	updated, err := mockRepo.Update(user)
	if err != nil {
		t.Fatalf("Update falló: %v", err)
	}

	if updated.Name != "Updated" {
		t.Errorf("Update() = %v, quiere Updated", updated.Name)
	}

	// Test Delete
	err = mockRepo.Delete(int(user.ID))
	if err != nil {
		t.Fatalf("Delete falló: %v", err)
	}

	_, err = mockRepo.FindByID(int(user.ID))
	if err == nil {
		t.Errorf("FindByID() después de Delete no falló, quiere error")
	}
}
