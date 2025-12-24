package workers

import (
	"context"
	"sync"
	"testing"
	"time"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWorkerPool(t *testing.T) {
	t.Run("creates worker pool with specified size", func(t *testing.T) {
		workerCount := 3
		queueSize := 10
		pool := NewWorkerPool(workerCount, queueSize)

		assert.NotNil(t, pool)
		assert.Equal(t, workerCount, pool.workerCount)
		assert.Equal(t, queueSize, cap(pool.tasks))
		assert.Equal(t, queueSize, cap(pool.results))
	})

	t.Run("creates worker pool with zero workers", func(t *testing.T) {
		pool := NewWorkerPool(0, 5)
		assert.NotNil(t, pool)
		assert.Equal(t, 0, pool.workerCount)
	})

	t.Run("creates worker pool with zero queue size", func(t *testing.T) {
		pool := NewWorkerPool(2, 0)
		assert.NotNil(t, pool)
		assert.Equal(t, 0, cap(pool.tasks))
	})
}

func TestWorkerPool_StartStop(t *testing.T) {
	t.Run("start and stop workers", func(t *testing.T) {
		pool := NewWorkerPool(2, 5)

		// Start workers
		pool.Start(context.Background())

		// Give workers time to start
		time.Sleep(10 * time.Millisecond)

		// Stop should not panic
		assert.NotPanics(t, func() {
			pool.Stop(context.Background())
		})
	})

	t.Run("stop without starting", func(t *testing.T) {
		pool := NewWorkerPool(2, 5)

		// Stop should work even if not started
		assert.NotPanics(t, func() {
			pool.Stop(context.Background())
		})
	})

	t.Run("multiple starts don't panic", func(t *testing.T) {
		pool := NewWorkerPool(2, 5)

		pool.Start(context.Background())
		// Second start should not panic (though it's not ideal usage)
		assert.NotPanics(t, func() {
			pool.Start(context.Background())
		})

		pool.Stop(context.Background())
	})
}

func TestWorkerPool_SubmitAndProcess(t *testing.T) {
	t.Run("submit and process single task", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		order := &entities.Order{
			ID:     1,
			Status: valueobjects.OrderStatus("pending"),
		}

		newStatus := valueobjects.OrderStatus("processing")
		pool.Submit(context.Background(), order, "updateStatus", &newStatus)

		// Wait for result
		select {
		case result := <-pool.GetResults(context.Background()):
			assert.Equal(t, 1, result.ID)
			assert.Equal(t, valueobjects.OrderStatus("processing"), result.Status)
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for result")
		}
	})

	t.Run("submit multiple tasks", func(t *testing.T) {
		pool := NewWorkerPool(2, 10)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		orders := []*entities.Order{
			{ID: 1, Status: valueobjects.OrderStatus("pending")},
			{ID: 2, Status: valueobjects.OrderStatus("pending")},
			{ID: 3, Status: valueobjects.OrderStatus("pending")},
		}

		newStatus := valueobjects.OrderStatus("processing")
		for _, order := range orders {
			pool.Submit(context.Background(), order, "updateStatus", &newStatus)
		}

		// Collect results
		results := make(map[int]bool)
		for range orders {
			select {
			case result := <-pool.GetResults(context.Background()):
				results[result.ID] = true
				assert.Equal(t, valueobjects.OrderStatus("processing"), result.Status)
			case <-time.After(1 * time.Second):
				t.Fatal("timeout waiting for results")
			}
		}

		assert.Len(t, results, 3)
		assert.True(t, results[1])
		assert.True(t, results[2])
		assert.True(t, results[3])
	})

	t.Run("submit more tasks than workers", func(t *testing.T) {
		pool := NewWorkerPool(2, 100)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		const taskCount = 50
		for i := 0; i < taskCount; i++ {
			order := &entities.Order{ID: i}
			pool.Submit(context.Background(), order, "validate", nil)
		}

		// Should process all tasks
		processedCount := 0
		timeout := time.After(2 * time.Second)

		for processedCount < taskCount {
			select {
			case <-pool.GetResults(context.Background()):
				processedCount++
			case <-timeout:
				t.Fatalf("timeout, processed %d of %d tasks", processedCount, taskCount)
			}
		}

		assert.Equal(t, taskCount, processedCount)
	})
}

func TestWorkerPool_TaskTypes(t *testing.T) {
	t.Run("updateStatus task", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		order := &entities.Order{
			ID:     1,
			Status: valueobjects.OrderStatus("pending"),
		}

		newStatus := valueobjects.OrderStatus("shipped")
		pool.Submit(context.Background(), order, "updateStatus", &newStatus)

		result := <-pool.GetResults(context.Background())
		assert.Equal(t, valueobjects.OrderStatus("shipped"), result.Status)
	})

	t.Run("validate task", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		order := &entities.Order{
			ID:     1,
			Status: valueobjects.OrderStatus("pending"),
		}

		pool.Submit(context.Background(), order, "validate", nil)

		result := <-pool.GetResults(context.Background())
		assert.Equal(t, valueobjects.OrderStatus(entities.StatusProcessing), result.Status)
	})

	t.Run("calculate task", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		order := &entities.Order{
			ID: 1,
			Items: []entities.OrderItem{
				{ProductID: 1, Quantity: 2, Price: 10.0},
				{ProductID: 2, Quantity: 1, Price: 5.0},
			},
		}

		pool.Submit(context.Background(), order, "calculate", nil)

		result := <-pool.GetResults(context.Background())
		assert.Equal(t, 25.0, result.Total) // (2*10) + (1*5)
	})

	t.Run("complete task", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		order := &entities.Order{
			ID:     1,
			Status: valueobjects.OrderStatus("pending"),
		}

		pool.Submit(context.Background(), order, "complete", nil)

		result := <-pool.GetResults(context.Background())
		assert.Equal(t, valueobjects.OrderStatus(entities.StatusCompleted), result.Status)
	})

	t.Run("unknown task type leaves order unchanged", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		originalStatus := valueobjects.OrderStatus("original")
		order := &entities.Order{
			ID:     1,
			Status: originalStatus,
		}

		pool.Submit(context.Background(), order, "unknownType", nil)

		result := <-pool.GetResults(context.Background())
		// Unknown task type should not modify order
		assert.Equal(t, originalStatus, result.Status)
	})
}

func TestWorkerPool_Concurrency(t *testing.T) {
	t.Run("multiple workers process concurrently", func(t *testing.T) {
		const workerCount = 4
		const taskCount = 8

		pool := NewWorkerPool(workerCount, taskCount)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		startTime := time.Now()

		// Submit tasks
		for i := 0; i < taskCount; i++ {
			order := &entities.Order{ID: i}
			pool.Submit(context.Background(), order, "calculate", nil)
		}

		// Wait for all results
		for i := 0; i < taskCount; i++ {
			<-pool.GetResults(context.Background())
		}

		elapsed := time.Since(startTime)

		// With concurrency, should be faster than sequential
		// (though the dummy CPU work is minimal)
		t.Logf("Processed %d tasks with %d workers in %v", taskCount, workerCount, elapsed)
	})

	t.Run("no data race with concurrent access", func(t *testing.T) {
		pool := NewWorkerPool(2, 10)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		var wg sync.WaitGroup
		errors := make(chan error, 100)

		// Concurrent submissions
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				order := &entities.Order{ID: id}
				pool.Submit(context.Background(), order, "validate", nil)
			}(i)
		}

		// Concurrent result reading
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-pool.GetResults(context.Background()):
					// Success
				case <-time.After(100 * time.Millisecond):
					errors <- assert.AnError
				}
			}()
		}

		wg.Wait()
		close(errors)

		// Should have no errors
		for err := range errors {
			t.Errorf("concurrent access error: %v", err)
		}
	})
}

func TestWorkerPool_CalculateTotal(t *testing.T) {
	t.Run("calculate total for order with items", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)

		order := &entities.Order{
			Items: []entities.OrderItem{
				{Quantity: 2, Price: 10.0},
				{Quantity: 3, Price: 5.0},
			},
		}

		pool.calculateTotal(order)
		assert.Equal(t, 35.0, order.Total) // (2*10) + (3*5)
	})

	t.Run("calculate total for order without items", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)

		order := &entities.Order{
			Items: []entities.OrderItem{},
		}

		pool.calculateTotal(order)
		assert.Equal(t, 0.0, order.Total)
	})

	t.Run("calculate total with decimal prices", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)

		order := &entities.Order{
			Items: []entities.OrderItem{
				{Quantity: 2, Price: 9.99},
				{Quantity: 1, Price: 4.50},
			},
		}

		pool.calculateTotal(order)
		// 2*9.99 + 1*4.50 = 19.98 + 4.50 = 24.48
		assert.InDelta(t, 24.48, order.Total, 0.000001)
	})
}

func TestWorkerPool_ResultsChannel(t *testing.T) {
	t.Run("results channel is read-only from outside", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)

		results := pool.GetResults(context.Background())

		// Should be a receive-only channel
		_, ok := interface{}(results).(<-chan *entities.Order)
		assert.True(t, ok, "GetResults should return a receive-only channel")

		// Should not be able to send to it
		// This would be a compile-time error if uncommented:
		// results <- &entities.Order{}
	})

	t.Run("drain results after stop", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())

		// Submit some tasks
		for i := 0; i < 3; i++ {
			order := &entities.Order{ID: i}
			pool.Submit(context.Background(), order, "validate", nil)
		}

		// Stop workers
		pool.Stop(context.Background())
		// Should be able to read remaining results
		results := 0
		for result := range pool.GetResults(context.Background()) {
			require.NotNil(t, result)
			results++
		}

		assert.Equal(t, 3, results)
	})
}

func TestWorkerPool_EdgeCases(t *testing.T) {
	t.Run("submit nil status for updateStatus task", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		order := &entities.Order{
			ID:     1,
			Status: valueobjects.OrderStatus("pending"),
		}

		// Submit with nil status for updateStatus task
		pool.Submit(context.Background(), order, "updateStatus", nil)

		select {
		case result := <-pool.GetResults(context.Background()):
			assert.NotNil(t, result)
			// Status should remain unchanged since nil status was provided
			assert.Equal(t, valueobjects.OrderStatus("pending"), result.Status)
		case <-time.After(500 * time.Millisecond):
			t.Error("expected result")
		}
	})

	t.Run("submit after stop", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		pool.Stop(context.Background())

		// Submit after stop should panic (channel closed)
		assert.Panics(t, func() {
			pool.Submit(context.Background(), &entities.Order{}, "validate", nil)
		})
	})

	t.Run("zero capacity queue", func(t *testing.T) {
		pool := NewWorkerPool(1, 0) // Zero buffer
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		pool.Submit(context.Background(), &entities.Order{ID: 1}, "validate", nil)

		select {
		case result := <-pool.GetResults(context.Background()):
			require.NotNil(t, result)
			assert.Equal(t, 1, result.ID)
		case <-time.After(500 * time.Millisecond):
			t.Fatal("timeout waiting for result")
		}

		// Primera tarea debería procesarse
		pool.Submit(context.Background(), &entities.Order{ID: 1}, "validate", nil)

		select {
		case result := <-pool.GetResults(context.Background()):
			require.NotNil(t, result)
			assert.Equal(t, 1, result.ID)
		case <-time.After(500 * time.Millisecond):
			t.Fatal("timeout waiting for result")
		}

		pool.Submit(context.Background(), &entities.Order{ID: 2}, "validate", nil)

		select {
		case result := <-pool.GetResults(context.Background()):
			require.NotNil(t, result)
			assert.Equal(t, 2, result.ID)
		case <-time.After(500 * time.Millisecond):
			t.Fatal("timeout waiting for second result")
		}
	})

	t.Run("unknown task type", func(t *testing.T) {
		pool := NewWorkerPool(1, 5)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		originalStatus := valueobjects.OrderStatus("original")
		order := &entities.Order{
			ID:     1,
			Status: originalStatus,
		}

		pool.Submit(context.Background(), order, "unknownType", nil)

		result := <-pool.GetResults(context.Background())
		// Unknown task type should not modify order
		assert.Equal(t, originalStatus, result.Status)
	})
}

func TestOrderTask_Struct(t *testing.T) {
	t.Run("order task creation", func(t *testing.T) {
		order := &entities.Order{ID: 1}
		status := valueobjects.OrderStatus("test")

		task := OrderTask{
			Order:  order,
			Type:   "update",
			Status: &status,
		}

		assert.Equal(t, order, task.Order)
		assert.Equal(t, "update", task.Type)
		assert.Equal(t, &status, task.Status)
	})

	t.Run("order task with nil status", func(t *testing.T) {
		task := OrderTask{
			Order:  &entities.Order{},
			Type:   "validate",
			Status: nil,
		}

		assert.Nil(t, task.Status)
	})
}

// Benchmark tests
func BenchmarkWorkerPool(b *testing.B) {
	b.Run("single worker", func(b *testing.B) {
		pool := NewWorkerPool(1, b.N)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			order := &entities.Order{ID: i}
			pool.Submit(context.Background(), order, "validate", nil)
		}

		// Drain results
		for i := 0; i < b.N; i++ {
			<-pool.GetResults(context.Background())
		}
	})

	b.Run("multiple workers", func(b *testing.B) {
		pool := NewWorkerPool(4, b.N)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			order := &entities.Order{ID: i}
			pool.Submit(context.Background(), order, "validate", nil)
		}

		// Drain results
		for i := 0; i < b.N; i++ {
			<-pool.GetResults(context.Background())
		}
	})
}

// Integration-style test
func TestWorkerPool_Integration(t *testing.T) {
	t.Run("process multiple order types correctly", func(t *testing.T) {
		pool := NewWorkerPool(2, 10)
		pool.Start(context.Background())
		defer pool.Stop(context.Background())

		testCases := []struct {
			name     string
			order    *entities.Order
			taskType string
			status   *valueobjects.OrderStatus
			validate func(t *testing.T, result *entities.Order)
		}{
			{
				name: "calculate total",
				order: &entities.Order{
					ID:    1,
					Items: []entities.OrderItem{{ProductID: 1, Quantity: 2, Price: 15.5}},
				},
				taskType: "calculate",
				validate: func(t *testing.T, result *entities.Order) {
					assert.Equal(t, 31.0, result.Total) // 2 * 15.5
				},
			},
			{
				name: "update status to processing",
				order: &entities.Order{
					ID:     2,
					Status: valueobjects.OrderStatus("pending"),
				},
				taskType: "updateStatus",
				status:   func() *valueobjects.OrderStatus { s := valueobjects.OrderStatus("processing"); return &s }(),
				validate: func(t *testing.T, result *entities.Order) {
					assert.Equal(t, valueobjects.OrderStatus("processing"), result.Status)
				},
			},
			{
				name: "complete order",
				order: &entities.Order{
					ID:     3,
					Status: valueobjects.OrderStatus("pending"),
				},
				taskType: "complete",
				validate: func(t *testing.T, result *entities.Order) {
					assert.Equal(t, valueobjects.OrderStatus(entities.StatusCompleted), result.Status)
				},
			},
		}

		// Enviar todas las tareas
		for _, tc := range testCases {
			pool.Submit(context.Background(), tc.order, tc.taskType, tc.status)
		}

		// Recibir y verificar resultados
		results := make(map[int]*entities.Order)
		for range testCases {
			select {
			case result := <-pool.GetResults(context.Background()):
				require.NotNil(t, result)
				results[result.ID] = result
			case <-time.After(1 * time.Second):
				t.Fatal("timeout waiting for results")
			}
		}

		// Ejecutar validaciones
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, ok := results[tc.order.ID]
				require.True(t, ok, "result for order %d not found", tc.order.ID)
				tc.validate(t, result)
			})
		}
	})
}
