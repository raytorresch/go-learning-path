package input

import (
	"context"
	"user-management/internal/domain/entities"

	"github.com/google/uuid"
)

type OrderService interface {
	PlaceOrder(ctx context.Context, userID uuid.UUID, items []entities.OrderItem) (*entities.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	GetAllOrders(ctx context.Context) ([]*entities.Order, error)
	CancelOrder(ctx context.Context, id uuid.UUID) error
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error
}
