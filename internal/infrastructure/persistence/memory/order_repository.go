package memory

import (
	"context"
	"sync"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/ports/output"

	"github.com/google/uuid"
)

type OrderRepository struct {
	mutex sync.RWMutex
}

var orders = map[int]*entities.Order{}

var _ output.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

// Delete implements [output.OrderRepository].
func (o *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	_, exists := orders[int(id.ID())]
	if !exists {
		return nil
	}

	delete(orders, int(id.ID()))
	return nil
}

// FindByID implements [output.OrderRepository].
func (o *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	order, exists := orders[int(id.ID())]
	if !exists {
		return nil, nil
	}

	return order, nil
}

// GetAllOrders implements [output.OrderRepository].
func (o *OrderRepository) GetAllOrders(ctx context.Context) ([]*entities.Order, error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	var allOrders []*entities.Order
	for _, order := range orders {
		allOrders = append(allOrders, order)
	}

	return allOrders, nil
}

// Save implements [output.OrderRepository].
func (o *OrderRepository) Save(ctx context.Context, order entities.Order) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	order.ID = int(uuid.New().ID())

	orders[int(order.ID)] = &order
	return nil
}

// Update implements [output.OrderRepository].
func (o *OrderRepository) Update(ctx context.Context, order *entities.Order) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	orders[int(order.ID)] = order
	return nil
}
