package main

import (
	"fmt"
	"log"
	"user-management/internal/models"
	"user-management/internal/processors"
	"user-management/internal/service"
	"user-management/internal/storage"
)

func main() {
	// init deps
	userRepo := storage.NewUserRepository()
	userService := service.NewUserService(userRepo)

	fmt.Println("=== Management Users System===")

	// create users
	user1, err := userService.CreateUser("Jhon", "jhon@mail.com", 30, true)
	if err != nil {
		log.Fatalf("Error creating user1: %v", err)
	}
	fmt.Println("Created User:", user1.String())

	user2, err := userService.CreateUser("Alice", "alice@mail.com", 25, true)
	if err != nil {
		log.Fatalf("Error creating user2: %v", err)
	}
	fmt.Println("Created User:", user2.String())

	// List all users
	fmt.Println("\n--- Users List ---")
	users, err := userService.GetAllUsers()
	if err != nil {
		log.Fatalf("Error getting users: %v", err)
	}
	for _, user := range users {
		fmt.Printf("%s\n", user.String())
	}

	// get user by Id
	fmt.Println("\n--- Get User by ID ---")
	user, err := userService.GetUserById(user1.ID)
	if err != nil {
		log.Fatalf("Error getting user by ID: %v", err)
	}
	fmt.Printf("User found: %s\n", user.String())

	// update user
	fmt.Println("\n--- Update User ---")
	updatedUser, err := userService.UpdateUser(user2.ID, "Alice Smith", "alice.smith@mail.com", 26, true)
	if err != nil {
		log.Fatalf("Error updating user: %v", err)
	}
	fmt.Printf("Updated User: %s\n", updatedUser.String())

	// deactivate user
	fmt.Println("\n--- Deactivate User ---")
	err = userService.DeactivateUser(user1.ID)
	if err != nil {
		log.Fatalf("Error deactivating user: %v", err)
	}
	deactivatedUser, err := userService.GetUserById(user1.ID)
	if err != nil {
		log.Fatalf("Error getting deactivated user: %v", err)
	}
	fmt.Printf("Deactivated User: %s\n", deactivatedUser.String())

	// user statistics
	fmt.Println("\n--- User Statistics ---")
	stats, err := userService.GetUserStatistics()
	if err != nil {
		log.Fatalf("Error getting user statistics: %v", err)
	}
	fmt.Printf("User Statistics: %+v\n", stats)

	// DEMO INTERFACES Y POLIMORFISMO
	fmt.Println("\n=== SISTEMA DE NOTIFICACIONES ===")

	notificationService := service.NewNotificationService()

	// Crear notificación
	notification := &models.Notification{
		UserID:  int(user1.ID),
		Title:   "Bienvenido",
		Message: "Gracias por registrarte",
		Type:    "system",
	}

	// Demostrar polimorfismo
	fmt.Printf("Canales disponibles: %v\n", notificationService.GetAvailableNotifiers())

	err = notificationService.SendToAllChannels(notification, user1)
	if err != nil {
		log.Printf("Error enviando notificación: %v", err)
	} else {
		fmt.Println("Notificaciones enviadas a todos los canales")
	}

	fmt.Println("\n=== SISTEMA DE PROCESAMIENTO DE TAREAS ===")

	// Inicializar servicios
	taskRepo := storage.NewTaskRepository()
	taskProcessor := processors.NewTaskProcessor()
	taskService := service.NewTaskService(taskRepo, taskProcessor)

	// Crear tareas variadas
	task1, _ := taskService.CreateAndProcess("email_report", int(user1.ID), 3)
	task2, _ := taskService.CreateAndProcess("data_cleanup", int(user2.ID), 5)
	task3, _ := taskService.CreateAndProcess("unknown_task", int(user1.ID), 1)

	// Procesamiento por lotes (variadic)
	completed, failed := taskService.BatchProcess(task1, task2, task3)
	fmt.Printf("Procesamiento por lotes: %d completadas, %d fallidas\n", completed, failed)

	// Filtrar tareas de alta prioridad
	highPriority := taskService.GetHighPriorityTasks(4)
	fmt.Printf("Tareas de alta prioridad (>=4): %d tareas\n", len(highPriority))

	// Demostrar switch con tipos (type switch)
	fmt.Println("\n Analizando tipos de datos en tareas:")
	for _, task := range []*models.Task{task1, task2, task3} {
		switch data := task.Data.(type) {
		case string:
			fmt.Printf("  Tarea %d: datos tipo string -> %s\n", task.ID, data)
		case nil:
			fmt.Printf("  Tarea %d: sin datos\n", task.ID)
		default:
			fmt.Printf("  Tarea %d: tipo desconocido: %T\n", task.ID, data)
		}
	}
}
