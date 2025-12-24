package output

import (
	"context"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/valueobjects"
)

// OrderWorker es un puerto para procesamiento asíncrono de órdenes
type OrderWorker interface {
	Start(ctx context.Context) error
	Submit(ctx context.Context, order *entities.Order, taskType string, status *valueobjects.OrderStatus) error
	GetResults(ctx context.Context) <-chan *entities.Order

	// Stop detiene el worker
	Stop(ctx context.Context) error
	IsStopped(ctx context.Context) bool
	WorkerCount(ctx context.Context) int
}

// OrderTask define una tarea para el worker
type OrderTask struct {
	Type   string
	Status *valueobjects.OrderStatus
}
