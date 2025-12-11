package main

import (
	"fmt"
	"time"
	"user-management/internal/domain/entities"
	"user-management/internal/infrastructure/services"
	"user-management/internal/infrastructure/storage"
)

func main() {
	// init deps
	userRepo := storage.NewMemoryUserRepository()
	taskRepo := storage.NewMemoryTaskRepository()

	userService := services.NewUserService(userRepo, taskRepo)

	user, err := userService.RegisterUser("Juan", "juan@test.com", 30)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Usuario creado: %s\n", user.Name)
	fmt.Printf("Email: %s\n", user.Email())

	// Demonstrar concurrencia
	demonstrateConcurrency()
}

func demonstrateConcurrency() {
	fmt.Println("\n=== DEMOSTRACIÓN DE CONCURRENCIA ===")

	// 1. Repository con channels
	orderRepo := storage.NewConcurrentOrderRepository(10)
	defer orderRepo.Close()

	// 2. Servicio con worker pool
	orderService := services.NewConcurrentOrderService(orderRepo, 3)

	// 3. Crear pedidos de prueba
	orders := createSampleOrders()

	// 4. Procesamiento batch concurrente
	fmt.Println("Procesando 5 pedidos concurrentemente...")
	start := time.Now()

	results := orderService.ProcessBatch(orders)

	elapsed := time.Since(start)
	fmt.Printf("✅ Tiempo total: %v (sería ~500ms secuencial)\n", elapsed)

	// 5. Mostrar resultados
	for _, order := range results {
		fmt.Printf("Pedido %d: %s - Total: $%.2f\n",
			order.ID, order.Status, order.Total)
	}

	// 6. Demostrar streaming con channels
	fmt.Println("\n=== STREAMING CON CHANNELS ===")
	orderChan := make(chan *entities.Order, 3)

	// Goroutine productora
	go func() {
		for i := 1; i <= 3; i++ {
			order := &entities.Order{
				UserID: i,
				Items: []entities.OrderItem{
					{ProductID: i, Name: fmt.Sprintf("Producto %d", i), Quantity: 2, Price: 25.99},
				},
				Status: entities.StatusReceived,
			}
			orderChan <- order
			time.Sleep(100 * time.Millisecond)
		}
		close(orderChan)
	}()

	// Consumir resultados
	resultChan := orderService.StreamOrders(orderChan)
	for order := range resultChan {
		fmt.Printf("📦 Pedido procesado: User %d - Total: $%.2f\n",
			order.UserID, order.Total)
	}
}

func createSampleOrders() []*entities.Order {
	orders := []*entities.Order{}
	for i := 1; i <= 5; i++ {
		order := &entities.Order{
			UserID: i,
			Items: []entities.OrderItem{
				{ProductID: i, Name: fmt.Sprintf("Producto %d", i), Quantity: 1 + i, Price: 10.0 * float64(i)},
			},
			Status: entities.StatusReceived,
		}
		orders = append(orders, order)
	}
	return orders
}
