package processors

import (
	"log"
	"sync"
	"user-management/internal/domain/entities"
)

type OrderTask struct {
	Order *entities.Order
	Type  string
}

type WorkerPool struct {
	tasks       chan OrderTask
	results     chan *entities.Order
	workerCount int
	wg          sync.WaitGroup
}

func NewWorkerPool(workerCount, queueSize int) *WorkerPool {
	return &WorkerPool{
		tasks:       make(chan OrderTask, queueSize),
		results:     make(chan *entities.Order, queueSize),
		workerCount: workerCount,
	}
}

// Start - Inicia el pool de workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		workerID := i + 1

		// Cada worker es una goroutine
		go wp.worker(workerID)
	}
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

		log.Printf("Worker %d procesando orden %d", id, task.Order.ID)

		// Procesar según tipo de tarea
		switch task.Type {
		case "validate":
			task.Order.Status = entities.OrderStatus(entities.StatusProcessing)
		case "calculate":
			wp.calculateTotal(task.Order)
		case "complete":
			task.Order.Status = entities.OrderStatus(entities.StatusCompleted)
		}

		// Enviar resultado
		wp.results <- task.Order
	}
}

func (wp *WorkerPool) calculateTotal(order *entities.Order) {
	total := 0.0
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total
}

// Submit - Envía tarea al pool
func (wp *WorkerPool) Submit(order *entities.Order, taskType string) {
	wp.tasks <- OrderTask{Order: order, Type: taskType}
}

// GetResults - Obtiene resultados (usando select para no bloquear)
func (wp *WorkerPool) GetResults() <-chan *entities.Order {
	return wp.results
}

// Stop - Detiene el pool de workers
func (wp *WorkerPool) Stop() {
	close(wp.tasks)
	wp.wg.Wait()
	close(wp.results)
}
