// mocks/order_repository_mock.go
package mocks

import (
	"context"
	"user-management/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// OrderRepositoryMock es un mock para output.OrderRepository
type OrderRepositoryMock struct {
	mock.Mock
}

// Delete implementa output.OrderRepository
func (m *OrderRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// FindByID implementa output.OrderRepository
func (m *OrderRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.Order), args.Error(1)
}

// GetAllOrders implementa output.OrderRepository
func (m *OrderRepositoryMock) GetAllOrders(ctx context.Context) ([]*entities.Order, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entities.Order), args.Error(1)
}

// Save implementa output.OrderRepository
func (m *OrderRepositoryMock) Save(ctx context.Context, order entities.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

// Update implementa output.OrderRepository
func (m *OrderRepositoryMock) Update(ctx context.Context, order *entities.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

// Métodos helper para facilitar la configuración de tests

// SetupDelete configura el mock para el método Delete
func (m *OrderRepositoryMock) SetupDelete(id uuid.UUID, err error) *mock.Call {
	return m.On("Delete", mock.Anything, id).Return(err)
}

// SetupFindByID configura el mock para el método FindByID
func (m *OrderRepositoryMock) SetupFindByID(id uuid.UUID, order *entities.Order, err error) *mock.Call {
	return m.On("FindByID", mock.Anything, id).Return(order, err)
}

// SetupGetAllOrders configura el mock para el método GetAllOrders
func (m *OrderRepositoryMock) SetupGetAllOrders(orders []*entities.Order, err error) *mock.Call {
	return m.On("GetAllOrders", mock.Anything).Return(orders, err)
}

// SetupSave configura el mock para el método Save
func (m *OrderRepositoryMock) SetupSave(order entities.Order, err error) *mock.Call {
	return m.On("Save", mock.Anything, order).Return(err)
}

// SetupUpdate configura el mock para el método Update
func (m *OrderRepositoryMock) SetupUpdate(order *entities.Order, err error) *mock.Call {
	return m.On("Update", mock.Anything, order).Return(err)
}

// SetupSaveWithIDGeneration configura Save para que actualice el ID
func (m *OrderRepositoryMock) SetupSaveWithIDGeneration(order entities.Order, generatedID int, err error) *mock.Call {
	return m.On("Save", mock.Anything, mock.MatchedBy(func(arg entities.Order) bool {
		// Actualiza el ID en el argumento recibido
		order.ID = generatedID
		return true
	})).Return(err)
}

// SetupSaveAny configura Save para aceptar cualquier orden
func (m *OrderRepositoryMock) SetupSaveAny(err error) *mock.Call {
	return m.On("Save", mock.Anything, mock.AnythingOfType("entities.Order")).Return(err)
}

// SetupFindByIDNotFound configura FindByID para devolver nil (no encontrado)
func (m *OrderRepositoryMock) SetupFindByIDNotFound(id uuid.UUID) *mock.Call {
	return m.On("FindByID", mock.Anything, id).Return(nil, nil)
}

// SetupFindByIDError configura FindByID para devolver un error
func (m *OrderRepositoryMock) SetupFindByIDError(id uuid.UUID, err error) *mock.Call {
	return m.On("FindByID", mock.Anything, id).Return(nil, err)
}

// VerifyExpectations verifica que todas las expectativas se cumplieron
func (m *OrderRepositoryMock) VerifyExpectations(t mock.TestingT) bool {
	return m.AssertExpectations(t)
}
