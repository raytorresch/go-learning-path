package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"user-management/internal/infrastructure/http/handlers"
	"user-management/internal/infrastructure/middlewares"
	"user-management/internal/infrastructure/persistence/memory"
	// "user-management/internal/infrastructure/storage"
)

func main() {
	// Configurar modo Gin
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar dependencias
	userRepo := memory.NewUserRepository()
	// orderRepo := storage.NewConcurrentOrderRepository(100)

	// userService := services.NewUserService(userRepo, nil)
	// orderService := services.NewConcurrentOrderService(orderRepo, 5)

	// Crear router
	router := gin.Default()

	// Middlewares globales
	router.Use(middlewares.LoggingMiddleware())
	router.Use(gin.Recovery()) // Recupera de panics

	// CORS básico
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rutas públicas
	public := router.Group("/api/v1")
	{
		healthHandler := handlers.NewHealthHandler()
		healthHandler.RegisterRoutes(public)
	}

	// Rutas protegidas (con autenticación)
	api := router.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware()) // Middleware de auth
	{
		// Users
		userHandler := handlers.NewUserHandler(userRepo)
		userHandler.RegisterRoutes(api)

		// Orders
		// orderHandler := handlers.NewOrderHandler(orderService)
		// orderHandler.RegisterRoutes(api)
	}

	// Servir documentación
	router.Static("/docs", "./docs")

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado en http://localhost:%s", port)
	log.Printf("Documentación en http://localhost:%s/docs", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
