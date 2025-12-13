package storage

import (
	"fmt"
	"sync"
	"time"
	"user-management/internal/domain/entities"
)

type ConcurrentOrderRepository struct {
	orders    map[int]*entities.Order
	mutex     sync.RWMutex
	nextID    int
	orderChan chan *entities.Order // Channel para pedidos entrantes
	done      chan struct{}        // Channel para señal de parada
}

func NewConcurrentOrderRepository(bufferSize int) *ConcurrentOrderRepository {
	repo := &ConcurrentOrderRepository{
		orders:    make(map[int]*entities.Order),
		nextID:    1,
		orderChan: make(chan *entities.Order, bufferSize), // Buffered channel
		done:      make(chan struct{}),
	}

	// Goroutine que procesa pedidos concurrentemente
	go repo.processOrders()

	return repo
}

// processOrders - Goroutine que consume del channel
func (r *ConcurrentOrderRepository) processOrders() {
	for {
		select {
		case order := <-r.orderChan:
			// Procesar pedido de forma concurrente
			r.saveOrder(order)
		case <-r.done:
			// Señal para terminar la goroutine
			return
		}
	}
}

// SaveAsync - Envía pedido al channel (no bloqueante si hay buffer)
func (r *ConcurrentOrderRepository) SaveAsync(order *entities.Order) {
	// Simular procesamiento asíncrono
	go func(o *entities.Order) {
		o.Status = entities.OrderStatus(entities.StatusProcessing)
		time.Sleep(time.Millisecond * 100) // Simular trabajo
		r.orderChan <- o
	}(order)
}

// saveOrder - Guarda el pedido (llamado por la goroutine)
func (r *ConcurrentOrderRepository) saveOrder(order *entities.Order) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	order.ID = r.nextID
	order.CreatedAt = time.Now()
	r.orders[order.ID] = order
	r.nextID++

	// Simular completado
	time.Sleep(time.Millisecond * 50)
	order.Status = entities.OrderStatus(entities.StatusCompleted)
	order.CompletedAt = time.Now()
}

// FindByID - Búsqueda síncrona
func (r *ConcurrentOrderRepository) FindByID(id int) (*entities.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

// Close - Limpieza de goroutines
func (r *ConcurrentOrderRepository) Close() {
	close(r.done)
	close(r.orderChan)
}
