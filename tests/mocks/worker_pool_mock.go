package mocks

import (
	"context"
	"sync"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/valueobjects"

	"github.com/stretchr/testify/mock"
)

// WorkerPoolMock implementa output.OrderWorker
type WorkerPoolMock struct {
	mock.Mock
	mu          sync.RWMutex
	submits     []*SubmitCall
	startCalled bool
	stopCalled  bool
}

// SubmitCall registra una llamada a Submit
type SubmitCall struct {
	Ctx      context.Context
	Order    *entities.Order
	TaskType string
	Status   *valueobjects.OrderStatus
}

// NewWorkerPoolMock crea un nuevo mock de WorkerPool
func NewWorkerPoolMock() *WorkerPoolMock {
	return &WorkerPoolMock{
		submits: make([]*SubmitCall, 0),
	}
}

// Start implementa output.OrderWorker.Start
func (m *WorkerPoolMock) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.startCalled = true
	args := m.Called(ctx)
	return args.Error(0)
}

// Submit implementa output.OrderWorker.Submit
func (m *WorkerPoolMock) Submit(ctx context.Context, order *entities.Order, taskType string, status *valueobjects.OrderStatus) error {
	m.mu.Lock()
	call := &SubmitCall{
		Ctx:      ctx,
		Order:    order,
		TaskType: taskType,
		Status:   status,
	}
	m.submits = append(m.submits, call)
	m.mu.Unlock()

	args := m.Called(ctx, order, taskType, status)
	return args.Error(0)
}

// GetResults implementa output.OrderWorker.GetResults
func (m *WorkerPoolMock) GetResults(ctx context.Context) <-chan *entities.Order {
	args := m.Called(ctx)

	// IMPORTANTE: Usar type assertion correcta para canal receive-only
	if args.Get(0) == nil {
		return nil
	}

	// Convertir chan *entities.Order a <-chan *entities.Order
	ch := args.Get(0).(chan *entities.Order)
	return ch
}

// Stop implementa output.OrderWorker.Stop
func (m *WorkerPoolMock) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stopCalled = true
	args := m.Called(ctx)
	return args.Error(0)
}

// IsStopped implementa output.OrderWorker.IsStopped
func (m *WorkerPoolMock) IsStopped(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// WorkerCount implementa output.OrderWorker.WorkerCount
func (m *WorkerPoolMock) WorkerCount(ctx context.Context) int {
	args := m.Called(ctx)
	return args.Int(0)
}

// Métodos helper simplificados:

// SetupSubmitSuccess configura Submit para éxito
func (m *WorkerPoolMock) SetupSubmitSuccess(order *entities.Order, taskType string, status *valueobjects.OrderStatus) *mock.Call {
	return m.On("Submit", mock.Anything, order, taskType, status).Return(nil)
}

// SetupGetResultsWithOrder configura GetResults para devolver una orden
func (m *WorkerPoolMock) SetupGetResultsWithOrder(order *entities.Order) *mock.Call {
	ch := make(chan *entities.Order, 1)
	ch <- order
	close(ch)
	return m.On("GetResults", mock.Anything).Return(ch)
}

// SetupGetResultsNil configura GetResults para devolver nil
func (m *WorkerPoolMock) SetupGetResultsNil() *mock.Call {
	ch := make(chan *entities.Order, 1)
	ch <- nil
	close(ch)
	return m.On("GetResults", mock.Anything).Return(ch)
}

// SetupGetResultsEmpty configura GetResults para devolver canal vacío
func (m *WorkerPoolMock) SetupGetResultsEmpty() *mock.Call {
	ch := make(chan *entities.Order)
	close(ch)
	return m.On("GetResults", mock.Anything).Return(ch)
}
