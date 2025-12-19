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

func (s *UserService) RegisterUser(name, email string, age int, password string) (*entities.User, error) {
	user, _ := entities.NewUser(name, email, age, password)
	return s.userRepo.Save(user)
}

func (s *UserService) GetUserByID(id int) (*entities.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetAllUsers() ([]*entities.User, error) {
	return s.userRepo.FindAll()
}
