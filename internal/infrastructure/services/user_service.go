package services

import (
	"user-management/internal/domain/entities"
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

func (s *UserService) RegisterUser(name, email string, age int) (*entities.User, error) {
	user := entities.NewUser(name, email, age, true)
	return s.userRepo.Save(user)
}
