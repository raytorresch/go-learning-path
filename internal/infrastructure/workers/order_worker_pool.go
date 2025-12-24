package workers

import (
	"context"
	"log"
	"sync"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/ports/output"
	"user-management/internal/domain/valueobjects"
)

type OrderTask struct {
	Order  *entities.Order
	Type   string
	Status *valueobjects.OrderStatus
}

type WorkerPool struct {
	tasks       chan OrderTask
	results     chan *entities.Order
	workerCount int
	wg          sync.WaitGroup
	mu          sync.RWMutex
	stopped     bool
}

var _ output.OrderWorker = (*WorkerPool)(nil)

func NewWorkerPool(workerCount, queueSize int) *WorkerPool {
	return &WorkerPool{
		tasks:       make(chan OrderTask, queueSize),
		results:     make(chan *entities.Order, queueSize),
		workerCount: workerCount,
	}
}

// Start - Inicia el pool de workers
func (wp *WorkerPool) Start(ctx context.Context) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.stopped {
		log.Println("WorkerPool ya está detenido, no se puede iniciar")
		return nil
	}

	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		workerID := i + 1

		go wp.worker(workerID)
	}
	log.Printf("WorkerPool iniciado con %d workers", wp.workerCount)
	return nil
}

// worker - Goroutine individual
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for task := range wp.tasks {
		// Trabajo CPU para ver paralelismo
		total := 0
		for i := range 1000000 {
			total += i * id
		}

		if task.Order == nil {
			log.Printf("Worker %d: Orden nil recibida para tipo %s", id, task.Type)

			wp.results <- nil
			continue
		}

		log.Printf("Worker %d procesando orden %d", id, task.Order.ID)

		switch task.Type {
		case "updateStatus":
			if task.Status != nil {
				task.Order.Status = *task.Status
			}
		case "validate":
			task.Order.Status = valueobjects.OrderStatus(entities.StatusProcessing)
		case "calculate":
			wp.calculateTotal(task.Order)
		case "complete":
			task.Order.Status = valueobjects.OrderStatus(entities.StatusCompleted)
		}

		wp.results <- task.Order
	}
}

func (wp *WorkerPool) calculateTotal(order *entities.Order) {
	if order == nil {
		return
	}

	total := 0.0

	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total
}

// Submit - Envía tarea al pool
func (wp *WorkerPool) Submit(ctx context.Context, order *entities.Order, taskType string, status *valueobjects.OrderStatus) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if wp.stopped {
		panic("cannot submit to stopped WorkerPool")
	}

	wp.tasks <- OrderTask{Order: order, Type: taskType, Status: status}

	return nil
}

// GetResults - Obtiene resultados (usando select para no bloquear)
func (wp *WorkerPool) GetResults(ctx context.Context) <-chan *entities.Order {
	return wp.results
}

// Stop - Detiene el pool de workers
func (wp *WorkerPool) Stop(ctx context.Context) error {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.stopped {
		return nil
	}

	wp.stopped = true
	close(wp.tasks)
	wp.wg.Wait()
	close(wp.results)
	log.Println("WorkerPool detenido")
	return nil
}

func (wp *WorkerPool) IsStopped(ctx context.Context) bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.stopped
}

func (wp *WorkerPool) WorkerCount(ctx context.Context) int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.workerCount
}
