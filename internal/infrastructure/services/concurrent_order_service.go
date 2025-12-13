package services

import (
	"sync"
	"time"
	"user-management/internal/domain/entities"
	"user-management/internal/infrastructure/processors"
	"user-management/internal/infrastructure/storage"
)

type ConcurrentOrderService struct {
	repo *storage.ConcurrentOrderRepository
	// processor *processors.OrderTask
	wp *processors.WorkerPool
}

func NewConcurrentOrderService(
	repo *storage.ConcurrentOrderRepository,
	workerCount int,
) *ConcurrentOrderService {

	// Crear worker pool
	wp := processors.NewWorkerPool(workerCount, 100)
	wp.Start()

	return &ConcurrentOrderService{
		repo: repo,
		wp:   wp,
	}
}

// ProcessBatch - Procesa múltiples pedidos concurrentemente
func (s *ConcurrentOrderService) ProcessBatch(orders []*entities.Order) []*entities.Order {
	var wg sync.WaitGroup
	results := make([]*entities.Order, len(orders))

	for i, order := range orders {
		wg.Add(1)

		// Cada pedido en su propia goroutine
		go func(idx int, o *entities.Order) {
			defer wg.Done()
			orderCopy := *o
			orderCopy.ID = idx + 1
			// Enviar al worker pool
			s.wp.Submit(&orderCopy, "calculate")

			// Usar select con timeout
			select {
			case result := <-s.wp.GetResults():
				results[idx] = result
			case <-time.After(2 * time.Second):
				orderCopy.Status = entities.OrderStatus(entities.StatusFailed)
				results[idx] = &orderCopy
			}
		}(i, order)
	}

	wg.Wait()
	return results
}

// StreamOrders - Procesamiento con pipeline de channels
func (s *ConcurrentOrderService) StreamOrders(orderChan <-chan *entities.Order) <-chan *entities.Order {
	resultChan := make(chan *entities.Order, 10)

	go func() {
		var wg sync.WaitGroup

		for order := range orderChan {
			wg.Add(1)

			go func(o *entities.Order) {
				defer wg.Done()

				// Pipeline: validate → calculate → complete
				s.wp.Submit(o, "validate")
				<-s.wp.GetResults()

				s.wp.Submit(o, "calculate")
				<-s.wp.GetResults()

				s.wp.Submit(o, "complete")
				completedOrder := <-s.wp.GetResults()

				resultChan <- completedOrder
			}(order)
		}

		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}
