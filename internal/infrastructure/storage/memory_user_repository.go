package storage

import (
	"fmt"
	"sync"
	"time"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/repositories"
)

// MemoryUserRepository implementa UserRepository
type MemoryUserRepository struct {
	users  map[int]*entities.User
	mutex  sync.RWMutex
	nextID int
}

// Garantiza que MemoryUserRepository cumple con la interfaz
var _ repositories.UserRepository = (*MemoryUserRepository)(nil)

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:  make(map[int]*entities.User),
		nextID: 1,
	}
}

func (r *MemoryUserRepository) Save(user *entities.User) (*entities.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Asignar ID si es nuevo
	if user.ID == 0 {
		user.ID = int64(r.nextID)
		r.nextID++
	}

	r.users[int(user.ID)] = user
	return user, nil
}

func (r *MemoryUserRepository) FindByID(id int) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[int(id)]
	if !exists {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return user, nil
}

func (r *MemoryUserRepository) FindAll() ([]*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *MemoryUserRepository) Update(user *entities.User) (*entities.User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.users[int(user.ID)]
	if !exists {
		return nil, fmt.Errorf("user with id %d not found", user.ID)
	}

	if err := user.Update(user.Name, user.Email(), user.Age, user.Active); err != nil {
		return nil, fmt.Errorf("fail validation: %w", err)
	}
	user.SetUpdatedAt(time.Now())
	return user, nil
}

func (r *MemoryUserRepository) Delete(id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.users[int(id)]
	if !exists {
		return fmt.Errorf("user with id %d not found", id)
	}

	delete(r.users, int(id))
	return nil
}
