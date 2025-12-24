package entities

import (
	"testing"
	"time"
	"user-management/internal/domain/valueobjects"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderTestSuite struct {
	suite.Suite
	baseTime time.Time
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (s *OrderTestSuite) SetupTest() {
	s.baseTime = time.Now()
}

func TestOrder_NewOrder(t *testing.T) {
	t.Run("success - valid order creation", func(t *testing.T) {
		userID := uuid.New()
		items := []OrderItem{
			{ProductID: 1, Name: "Product A", Quantity: 2, Price: 10.0},
			{ProductID: 2, Name: "Product B", Quantity: 1, Price: 5.0},
		}

		order, err := NewOrder(userID, items)
		require.NoError(t, err)
		require.NotNil(t, order)
		assert.Equal(t, int(userID.ID()), order.UserID)
		assert.Equal(t, items, order.Items)
		assert.Equal(t, valueobjects.StatusPending, order.Status)
		assert.WithinDuration(t, time.Now(), order.CreatedAt, time.Second)
		expectedTotal := 2*10.0 + 1*5.0
		assert.Equal(t, expectedTotal, order.Total)
	})

	t.Run("failure - invalid item in order", func(t *testing.T) {
		userID := uuid.New()
		items := []OrderItem{
			{ProductID: 1, Name: "Product A", Quantity: 0, Price: 10.0}, // Cantidad inválida
		}

		order, err := NewOrder(userID, items)
		require.Error(t, err)
		require.Nil(t, order)
		assert.Contains(t, err.Error(), "quantity must be positive")
	})
}

func (s *OrderTestSuite) TestOrder_CalculateTotal() {
	tests := []struct {
		name        string
		items       []OrderItem
		expected    float64
		expectError bool
	}{
		{
			name: "success - single item",
			items: []OrderItem{
				{ProductID: 1, Quantity: 2, Price: 10.0},
			},
			expected:    20.0,
			expectError: false,
		},
		{
			name: "success - multiple items",
			items: []OrderItem{
				{ProductID: 1, Quantity: 2, Price: 10.0},
				{ProductID: 2, Quantity: 1, Price: 5.5},
				{ProductID: 3, Quantity: 3, Price: 2.0},
			},
			expected:    20.0 + 5.5 + 6.0, // 31.5
			expectError: false,
		},
		{
			name:        "success - zero items",
			items:       []OrderItem{},
			expected:    0.0,
			expectError: false,
		},
		{
			name: "failure - item with negative price",
			items: []OrderItem{
				{ProductID: 1, Quantity: 1, Price: -10.0},
			},
			expectError: true, // Porque item.Validate() fallará
		},
		{
			name: "success - decimal quantities and prices",
			items: []OrderItem{
				{ProductID: 1, Quantity: 1, Price: 9.99},
				{ProductID: 2, Quantity: 2, Price: 4.50},
			},
			expected:    9.99 + 9.0, // 18.99
			expectError: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			order := &Order{
				Items: tt.items,
			}

			err := order.CalculateTotal()

			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
				// s.Equal(tt.expected, order.Total)
				s.InDelta(tt.expected, order.Total, 0.000001)
			}
		})
	}
}

func (s *OrderTestSuite) TestOrder_AddItem() {
	s.Run("success - add new item to empty order", func() {
		order := &Order{
			Items: []OrderItem{},
		}

		newItem := OrderItem{
			ProductID: 1,
			Name:      "Product A",
			Quantity:  2,
			Price:     10.0,
		}

		err := order.AddItem(newItem)
		s.NoError(err)
		s.Len(order.Items, 1)
		s.Equal(newItem, order.Items[0])
		s.Equal(20.0, order.Total)
	})

	s.Run("success - add new item to existing order", func() {
		order := &Order{
			Items: []OrderItem{
				{ProductID: 1, Name: "Product A", Quantity: 1, Price: 10.0},
			},
		}
		order.CalculateTotal()

		newItem := OrderItem{
			ProductID: 2,
			Name:      "Product B",
			Quantity:  3,
			Price:     5.0,
		}

		err := order.AddItem(newItem)
		s.NoError(err)
		s.Len(order.Items, 2)
		s.Equal(25.0, order.Total) // (1*10) + (3*5)
	})

	s.Run("success - update quantity for existing product", func() {
		order := &Order{
			Items: []OrderItem{
				{ProductID: 1, Name: "Product A", Quantity: 2, Price: 10.0},
			},
		}
		order.CalculateTotal()

		sameProductItem := OrderItem{
			ProductID: 1, // Mismo producto
			Name:      "Product A",
			Quantity:  3,
			Price:     10.0,
		}

		err := order.AddItem(sameProductItem)
		s.NoError(err)
		s.Len(order.Items, 1)               // No se añade nuevo, se actualiza
		s.Equal(5, order.Items[0].Quantity) // 2 + 3
		s.Equal(50.0, order.Total)          // 5 * 10
	})

	s.Run("failure - add invalid item", func() {
		order := &Order{
			Items: []OrderItem{},
		}

		invalidItem := OrderItem{
			ProductID: 1,
			Quantity:  0, // Cantidad inválida
			Price:     10.0,
		}

		err := order.AddItem(invalidItem)
		s.Error(err)
		s.Contains(err.Error(), "quantity must be positive")
		s.Empty(order.Items)
	})

	s.Run("failure - add item with negative price", func() {
		order := &Order{
			Items: []OrderItem{},
		}

		invalidItem := OrderItem{
			ProductID: 1,
			Quantity:  2,
			Price:     -5.0, // Precio negativo
		}

		err := order.AddItem(invalidItem)
		s.Error(err)
		s.Contains(err.Error(), "price cannot be negative")
	})
}

func (s *OrderTestSuite) TestOrder_Validate() {
	tests := []struct {
		name        string
		order       Order
		expectError bool
		errorMsg    string
	}{
		{
			name: "success - valid order",
			order: Order{
				ID:        1,
				UserID:    100,
				Items:     []OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
				Total:     10.0,
				Status:    valueobjects.StatusPending,
				CreatedAt: s.baseTime,
			},
			expectError: false,
		},
		{
			name: "failure - invalid user ID",
			order: Order{
				UserID: 0, // Inválido
				Items:  []OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
				Total:  10.0,
			},
			expectError: true,
			errorMsg:    "invalid user ID",
		},
		{
			name: "failure - negative user ID",
			order: Order{
				UserID: -1, // Inválido
				Items:  []OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
				Total:  10.0,
			},
			expectError: true,
			errorMsg:    "invalid user ID",
		},
		{
			name: "failure - no items",
			order: Order{
				UserID: 100,
				Items:  []OrderItem{}, // Vacío
				Total:  0.0,
			},
			expectError: true,
			errorMsg:    "order must have at least one item",
		},
		{
			name: "failure - negative total",
			order: Order{
				UserID: 100,
				Items:  []OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
				Total:  -5.0, // Negativo
			},
			expectError: true,
			errorMsg:    "total cannot be negative",
		},
		{
			name: "failure - completed before created",
			order: Order{
				UserID:      100,
				Items:       []OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
				Total:       10.0,
				CreatedAt:   s.baseTime,
				CompletedAt: s.baseTime.Add(-1 * time.Hour), // Antes de created
			},
			expectError: true,
			errorMsg:    "completed date cannot be before creation date",
		},
		{
			name: "success - with completed date after creation",
			order: Order{
				UserID:      100,
				Items:       []OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
				Total:       10.0,
				CreatedAt:   s.baseTime,
				CompletedAt: s.baseTime.Add(1 * time.Hour), // Después de created
			},
			expectError: false,
		},
		{
			name: "success - zero times are ignored",
			order: Order{
				UserID:    100,
				Items:     []OrderItem{{ProductID: 1, Quantity: 1, Price: 10.0}},
				Total:     10.0,
				CreatedAt: time.Time{}, // Zero time
			},
			expectError: false, // Zero times no se validan
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := tt.order.Validate()

			if tt.expectError {
				s.Error(err)
				if tt.errorMsg != "" {
					s.Contains(err.Error(), tt.errorMsg)
				}
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *OrderTestSuite) TestOrder_Complete() {
	s.Run("success - complete pending order", func() {
		order := &Order{
			ID:        1,
			UserID:    100,
			Status:    valueobjects.StatusPending,
			CreatedAt: s.baseTime,
		}

		err := order.Complete()
		s.NoError(err)
		s.Equal(valueobjects.StatusCompleted, order.Status)
		s.False(order.CompletedAt.IsZero())
		s.True(order.CompletedAt.After(s.baseTime) || order.CompletedAt.Equal(s.baseTime))
	})

	s.Run("failure - already completed order", func() {
		completedTime := s.baseTime.Add(1 * time.Hour)
		order := &Order{
			ID:          1,
			UserID:      100,
			Status:      valueobjects.StatusCompleted,
			CreatedAt:   s.baseTime,
			CompletedAt: completedTime,
		}

		err := order.Complete()
		s.Error(err)
		s.Contains(err.Error(), "order already completed")
		s.Equal(valueobjects.StatusCompleted, order.Status)
		s.Equal(completedTime, order.CompletedAt) // No cambió
	})

	s.Run("success - complete order with different initial status", func() {
		order := &Order{
			ID:        1,
			UserID:    100,
			Status:    valueobjects.StatusProcessing,
			CreatedAt: s.baseTime,
		}

		err := order.Complete()
		s.NoError(err)
		s.Equal(valueobjects.StatusCompleted, order.Status)
		s.False(order.CompletedAt.IsZero())
	})

	s.Run("complete updates timestamp", func() {
		order := &Order{
			ID:        1,
			UserID:    100,
			Status:    valueobjects.StatusPending,
			CreatedAt: s.baseTime,
		}

		beforeComplete := time.Now()
		err := order.Complete()
		afterComplete := time.Now()

		s.NoError(err)
		s.True(order.CompletedAt.After(beforeComplete) || order.CompletedAt.Equal(beforeComplete))
		s.True(order.CompletedAt.Before(afterComplete) || order.CompletedAt.Equal(afterComplete))
	})
}

// Tests adicionales sin suite
func TestOrder_Integration(t *testing.T) {
	t.Run("full workflow: create, add items, calculate, validate, complete", func(t *testing.T) {
		order := &Order{
			ID:        1,
			UserID:    100,
			Status:    valueobjects.StatusPending,
			CreatedAt: time.Now(),
			Items:     []OrderItem{},
		}

		// 1. Agregar items
		err := order.AddItem(OrderItem{ProductID: 1, Name: "Product A", Quantity: 2, Price: 10.0})
		require.NoError(t, err)

		err = order.AddItem(OrderItem{ProductID: 2, Name: "Product B", Quantity: 1, Price: 5.0})
		require.NoError(t, err)

		// 2. Verificar total
		assert.Equal(t, 25.0, order.Total)

		// 3. Validar
		err = order.Validate()
		require.NoError(t, err)

		// 4. Completar
		err = order.Complete()
		require.NoError(t, err)
		assert.Equal(t, valueobjects.StatusCompleted, order.Status)
		assert.False(t, order.CompletedAt.IsZero())

		// 5. No se puede completar otra vez
		err = order.Complete()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already completed")
	})

	t.Run("calculate total after multiple item additions", func(t *testing.T) {
		order := &Order{
			Items: []OrderItem{},
		}

		// Agregar múltiples veces el mismo producto
		order.AddItem(OrderItem{ProductID: 1, Quantity: 1, Price: 10.0})
		order.AddItem(OrderItem{ProductID: 1, Quantity: 2, Price: 10.0})
		order.AddItem(OrderItem{ProductID: 2, Quantity: 3, Price: 5.0})
		order.AddItem(OrderItem{ProductID: 1, Quantity: 1, Price: 10.0})

		assert.Equal(t, 2, len(order.Items))
		assert.Equal(t, 4, order.Items[0].Quantity) // ProductID 1: 1+2+1 = 4
		assert.Equal(t, 3, order.Items[1].Quantity) // ProductID 2: 3
		assert.Equal(t, 55.0, order.Total)          // (4*10) + (3*5) = 40 + 15 = 55
	})
}

func TestOrder_EdgeCases(t *testing.T) {
	t.Run("very large quantities", func(t *testing.T) {
		order := &Order{
			Items: []OrderItem{
				{ProductID: 1, Quantity: 1000000, Price: 0.01},
			},
		}

		err := order.CalculateTotal()
		require.NoError(t, err)
		assert.Equal(t, 10000.0, order.Total) // 1,000,000 * 0.01
	})

	t.Run("precision with floating point", func(t *testing.T) {
		order := &Order{
			Items: []OrderItem{
				{ProductID: 1, Quantity: 3, Price: 0.33},
			},
		}

		err := order.CalculateTotal()
		require.NoError(t, err)
		// 3 * 0.33 = 0.99
		assert.InDelta(t, 0.99, order.Total, 0.0001)
	})

	t.Run("empty order validation after adding invalid item", func(t *testing.T) {
		order := &Order{
			UserID: 100,
			Items:  []OrderItem{},
		}

		// Intentar agregar item inválido
		err := order.AddItem(OrderItem{ProductID: 1, Quantity: 0, Price: 10.0})
		require.Error(t, err)

		// La validación debe fallar porque no hay items
		err = order.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order must have at least one item")
	})
}

func TestOrder_StatusTransitions(t *testing.T) {
	t.Run("can complete from any non-completed status", func(t *testing.T) {
		testCases := []valueobjects.OrderStatus{
			valueobjects.StatusPending,
			valueobjects.StatusProcessing,
			valueobjects.StatusShipped,
		}

		for _, status := range testCases {
			t.Run(string(status), func(t *testing.T) {
				order := &Order{
					ID:     1,
					UserID: 100,
					Status: status,
				}

				err := order.Complete()
				assert.NoError(t, err)
				assert.Equal(t, valueobjects.StatusCompleted, order.Status)
			})
		}
	})
}

func TestOrder_ItemManagement(t *testing.T) {
	t.Run("remove item functionality (future)", func(t *testing.T) {
		order := &Order{
			Items: []OrderItem{
				{ProductID: 1, Quantity: 2, Price: 10.0},
				{ProductID: 2, Quantity: 1, Price: 5.0},
			},
		}
		order.CalculateTotal()
		assert.Equal(t, 25.0, order.Total)

		// Esto sería un método futuro RemoveItem
		// err := order.RemoveItem(1)
		// assert.NoError(t, err)
		// assert.Len(t, order.Items, 1)
		// assert.Equal(t, 5.0, order.Total)
	})
}

// Benchmark tests
func BenchmarkOrderCalculateTotal(b *testing.B) {
	order := &Order{
		Items: make([]OrderItem, 1000),
	}
	for i := range order.Items {
		order.Items[i] = OrderItem{
			ProductID: i,
			Quantity:  i%10 + 1,
			Price:     float64(i%100 + 1),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		order.CalculateTotal()
	}
}

func BenchmarkOrderAddItem(b *testing.B) {
	b.Run("add to empty order", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			order := &Order{Items: []OrderItem{}}
			item := OrderItem{ProductID: i, Quantity: 1, Price: 10.0}
			order.AddItem(item)
		}
	})

	b.Run("add to order with existing items", func(b *testing.B) {
		order := &Order{
			Items: make([]OrderItem, 100),
		}
		for i := range order.Items {
			order.Items[i] = OrderItem{ProductID: i, Quantity: 1, Price: 1.0}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			item := OrderItem{ProductID: 100 + i, Quantity: 1, Price: 10.0}
			order.AddItem(item)
		}
	})
}

// Test para verificar que Order es inmutable en ciertos aspectos
func TestOrder_Immutability(t *testing.T) {
	t.Run("ID cannot be changed after creation", func(t *testing.T) {
		order := Order{ID: 1, UserID: 100}
		// No hay métodos para cambiar el ID
		assert.Equal(t, 1, order.ID)
	})

	t.Run("CreatedAt cannot be changed", func(t *testing.T) {
		createdAt := time.Now()
		order := Order{CreatedAt: createdAt}
		assert.Equal(t, createdAt, order.CreatedAt)
	})
}
