package output

import (
	"context"
	"user-management/internal/domain/entities"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Save(ctx context.Context, order entities.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAllOrders(ctx context.Context) ([]*entities.Order, error)
}
