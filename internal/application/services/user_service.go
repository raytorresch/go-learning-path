package services

import (
	"context"
	"errors"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/ports/input"
	"user-management/internal/domain/ports/output"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserService struct {
	repo output.UserRepository
}

var _ input.UserService = (*UserService)(nil)

func NewUserService(repo output.UserRepository) input.UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, name string, email string, age int, password string) (*entities.User, error) {
	// 1. Crear usuario en dominio
	user, err := entities.NewUser(name, email, age, password)
	if err != nil {
		return nil, err
	}

	// 2. Validar unicidad de email (regla de negocio)
	existing, err := s.repo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// 3. Persistir
	if err := s.repo.Save(ctx, *user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, id uuid.UUID) (*entities.User, error) {

	return s.repo.FindByID(ctx, id)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*entities.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UserService) UpdateProfile(ctx context.Context, user *entities.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, uid)
}
