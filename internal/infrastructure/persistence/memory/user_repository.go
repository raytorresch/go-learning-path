package memory

import (
	"context"
	"fmt"
	"sync"
	"time"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/ports/output"

	"github.com/google/uuid"
)

// UserRepository implementa output
type UserRepository struct {
	mutex sync.RWMutex
}

var users = map[int]*entities.User{}

// Garantiza que UserRepository cumple con la interfaz
var _ output.UserRepository = (*UserRepository)(nil)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create implements output.UserPort.
func (u *UserRepository) Save(ctx context.Context, user entities.User) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	user.ID = uuid.New()

	users[int(user.ID.ID())] = &user
	return nil
}

// Delete implements output.UserPort.
func (u *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	_, exists := users[int(id.ID())]
	if !exists {
		return fmt.Errorf("user with id %s not found", id.String())
	}

	delete(users, int(id.ID()))
	return nil
}

// FindByEmail implements output.UserPort.
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	// search user by email
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user with email %s not found", email)
}

// FindByID implements output.UserPort.
func (u *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	user, exists := users[int(id.ID())]
	if !exists {
		return nil, fmt.Errorf("user with id %s not found", id.String())
	}

	return user, nil
}

// Update implements output.UserPort.
func (u *UserRepository) Update(ctx context.Context, user *entities.User) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	_, exists := users[int(user.ID.ID())]
	if !exists {
		return fmt.Errorf("user with id %s not found", user.ID.String())
	}

	if err := user.Update(user.Name, user.Email, user.Age, user.Active); err != nil {
		return fmt.Errorf("fail validation: %w", err)
	}
	user.SetUpdatedAt(time.Now())
	return nil
}

// GetAllUsers implements output.UserPort.
func (u *UserRepository) GetAllUsers(ctx context.Context) ([]*entities.User, error) {
	u.mutex.RLock()
	defer u.mutex.RUnlock()

	var allUsers []*entities.User
	for _, user := range users {
		allUsers = append(allUsers, user)
	}

	return allUsers, nil
}
