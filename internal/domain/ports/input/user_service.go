package input

import (
	"context"
	"user-management/internal/domain/entities"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, name string, email string, age int, password string) (*entities.User, error)
	GetUserProfile(ctx context.Context, id uuid.UUID) (*entities.User, error)
	UpdateProfile(ctx context.Context, user *entities.User) error
	GetAllUsers(ctx context.Context) ([]*entities.User, error)
}
