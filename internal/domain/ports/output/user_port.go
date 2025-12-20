package output

import (
	"context"
	"user-management/internal/domain/entities"

	"github.com/google/uuid"
)

// PORT (interfaz que define el contrato)
type UserPort interface {
	Create(ctx context.Context, user entities.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]*entities.User, error)
}
