package service

import (
	"user-management/internal/models"
	"user-management/internal/storage"
)

type UserService struct {
	repo *storage.UserRepository
}

func NewUserService(repo *storage.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(name, email string, age int, active bool) (*models.User, error) {
	user := &models.User{
		Name:   name,
		Email:  email,
		Age:    age,
		Active: active,
	}

	return s.repo.Create(user)
}

func (s *UserService) GetUserById(id int64) (*models.User, error) {
	return s.repo.FindById(id)
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) UpdateUser(id int64, name, email string, age int, active bool) (*models.User, error) {
	return s.repo.Update(id, name, email, age, active)
}

func (s *UserService) DeactivateUser(id int64) error {
	user, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	_, err = s.repo.Update(id, user.Name, user.Email, user.Age, false)
	return err
}

func (s *UserService) GetUserStatistics() (map[string]interface{}, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	totalUsers := len(users)
	if totalUsers == 0 {
		return map[string]interface{}{
			"total_users":    0,
			"active_users":   0,
			"inactive_users": 0,
			"average_age":    0,
		}, nil
	}

	activeUsers := 0
	totalAge := 0
	for _, user := range users {
		if user.Active {
			activeUsers++
		}
		totalAge += user.Age
	}

	avgAge := float64(totalAge) / float64(totalUsers)

	return map[string]interface{}{
		"total_users":    totalUsers,
		"active_users":   activeUsers,
		"inactive_users": totalUsers - activeUsers,
		"average_age":    avgAge,
	}, nil
}
