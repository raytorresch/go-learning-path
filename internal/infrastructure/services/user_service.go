package services

import (
	"user-management/internal/domain/models"
	"user-management/internal/domain/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
	taskRepo repositories.TaskRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	taskRepo repositories.TaskRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}

func (s *UserService) RegisterUser(name, email string, age int) (*models.User, error) {
	user := models.NewUser(name, email, age, true)
	return s.userRepo.Save(user)
}
