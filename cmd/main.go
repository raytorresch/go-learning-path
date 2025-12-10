package main

import (
	"fmt"
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
}
