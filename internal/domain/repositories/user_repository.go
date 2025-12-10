package repositories

import (
	"user-management/internal/domain/entities"
)

type UserRepository interface {
	Save(user *entities.User) (*entities.User, error)
	FindByID(id int) (*entities.User, error)
	FindAll() ([]*entities.User, error)
	Update(user *entities.User) (*entities.User, error)
	Delete(id int) error
}

// type UserRepository struct {
// 	users  map[int]*entities.User
// 	mutex  sync.RWMutex
// 	nextID int
// }

// func NewUserRepository() *UserRepository {
// 	return &UserRepository{
// 		users:  make(map[int]*entities.User),
// 		nextID: 1,
// 	}
// }

// func (r *UserRepository) Create(user *entities.User) (*entities.User, error) {
// 	if err := user.Validate(); err != nil {
// 		return nil, fmt.Errorf("fail validation: %w", err)
// 	}

// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	user.ID = int64(r.nextID)
// 	user.CreatedAt = time.Now()
// 	user.UpdatedAt = time.Now()

// 	r.users[r.nextID] = user
// 	r.nextID++

// 	return user, nil
// }

// func (r *UserRepository) FindById(id int64) (*entities.User, error) {
// 	r.mutex.RLock()
// 	defer r.mutex.RUnlock()

// 	user, exists := r.users[int(id)]
// 	if !exists {
// 		return nil, fmt.Errorf("user with id %d not found", id)
// 	}

// 	return user, nil
// }

// func (r *UserRepository) FindAll() ([]*entities.User, error) {
// 	r.mutex.RLock()
// 	defer r.mutex.RUnlock()

// 	users := make([]*entities.User, 0, len(r.users))
// 	for _, user := range r.users {
// 		users = append(users, user)
// 	}

// 	return users, nil
// }

// func (r *UserRepository) Update(id int64, name, email string, age int, active bool) (*entities.User, error) {
// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	user, exists := r.users[int(id)]
// 	if !exists {
// 		return nil, fmt.Errorf("user with id %d not found", id)
// 	}

// 	if err := user.Update(name, email, age, active); err != nil {
// 		return nil, fmt.Errorf("fail validation: %w", err)
// 	}

// 	user.UpdatedAt = time.Now()
// 	return user, nil
// }

// func (r *UserRepository) Delete(id int64) error {
// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	_, exists := r.users[int(id)]
// 	if !exists {
// 		return fmt.Errorf("user with id %d not found", id)
// 	}

// 	delete(r.users, int(id))
// 	return nil
// }
