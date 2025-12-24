package services

import (
	"context"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/ports/input"
	"user-management/internal/domain/ports/output"
	"user-management/internal/domain/valueobjects"

	"github.com/google/uuid"
)

type OrderService struct {
	repo   output.OrderRepository
	worker output.OrderWorker
}

var _ input.OrderService = (*OrderService)(nil)

func NewOrderService(repo output.OrderRepository, worker output.OrderWorker) input.OrderService {
	worker.Start(context.Background())
	return &OrderService{repo: repo, worker: worker}
}

// CancelOrder implements [input.OrderService].
func (o *OrderService) CancelOrder(ctx context.Context, id uuid.UUID) error {
	order, err := o.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return nil
	}

	order.Status = valueobjects.OrderStatus(entities.StatusCancelled)
	if err := o.repo.Update(ctx, order); err != nil {
		return err
	}
	return nil
}

// GetAllOrders implements [input.OrderService].
func (o *OrderService) GetAllOrders(ctx context.Context) ([]*entities.Order, error) {
	orders, err := o.repo.GetAllOrders(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderByID implements [input.OrderService].
func (o *OrderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*entities.Order, error) {
	if orderID == uuid.Nil {
		return nil, nil
	}

	order, err := o.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// PlaceOrder implements [input.OrderService].
func (o *OrderService) PlaceOrder(ctx context.Context, userID uuid.UUID, items []entities.OrderItem) (*entities.Order, error) {
	order, err := entities.NewOrder(userID, items)
	if err != nil {
		return nil, err
	}

	if err := o.repo.Save(ctx, *order); err != nil {
		return nil, err
	}

	return order, nil
}

// UpdateOrderStatus implements [input.OrderService].
func (o *OrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	order, err := o.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if order == nil {
		return nil
	}

	statusVO := valueobjects.OrderStatus(status)
	if err := o.worker.Submit(ctx, order, "updateStatus", &statusVO); err != nil {
		return err
	}

	resultsChan := o.worker.GetResults(ctx)

	select {
	case updatedOrder, ok := <-resultsChan:
		if !ok {
			// Canal cerrado sin resultados
			return nil
		}
		if updatedOrder == nil {
			// Resultado nil
			return nil
		}
		return o.repo.Update(ctx, updatedOrder)
	case <-ctx.Done():
		// Contexto cancelado o timeout
		return ctx.Err()
	}
}
