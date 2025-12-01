package main

import (
	"fmt"
	"log"
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
}
