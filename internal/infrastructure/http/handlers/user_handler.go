package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"user-management/internal/domain/entities"
	"user-management/internal/domain/ports/output"
)

type UserHandler struct {
	userPort output.UserPort
	validate *validator.Validate
}

func NewUserHandler(userPort output.UserPort) *UserHandler {
	return &UserHandler{
		userPort: userPort,
		validate: validator.New(),
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	// CRUD básico
	router.GET("/users", h.GetAllUsers)
	router.GET("/users/:id", h.GetUserByID)
	router.POST("/users", h.CreateUser)
	router.PUT("/users/:id", h.UpdateUser)
	router.DELETE("/users/:id", h.DeleteUser)

	// Rutas específicas
	router.GET("/users/:id/orders", h.GetUserOrders)
	router.POST("/users/:id/notify", h.NotifyUser)
}

// CreateUser demuestra binding y validación
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		Age      int    `json:"age" binding:"gte=0,lte=120"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// Binding automático con Gin
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Validación adicional
	if err := h.validate.Struct(req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user := h.userPort.Create(c.Request.Context(), req.Name, req.Email, req.Age, req.Password)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
		Message: "User created successfully",
	})
}

// GetUserByID con parámetros de ruta
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.userPort.FindByID(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, err)
		return
	}

	SuccessResponse(c, user)
}

// GetAllUsers con query parameters
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Ejemplo: /users?active=true&limit=10&offset=0
	// activeStr := c.DefaultQuery("active", "")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	// Filtrar usuarios (simplificado)
	users, _ := h.userPort.GetAllUsers(c.Request.Context())

	// Paginación básica
	start := offset
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	if start >= len(users) {
		SuccessResponse(c, []entities.User{})
		return
	}

	SuccessResponse(c, users[start:end])
}

// UpdateUser con actualización parcial (PATCH-like)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Aquí iría la lógica de actualización
	// Por ahora solo retornamos éxito
	SuccessResponse(c, gin.H{
		"id":      id,
		"updates": updates,
		"message": "User updated",
	})
}

// DeleteUser elimina un usuario por ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Aquí iría la lógica de eliminación
	// Por ahora solo retornamos éxito
	SuccessResponse(c, gin.H{
		"id":      id,
		"message": "User deleted",
	})
}

// GetUserOrders muestra nested resources
func (h *UserHandler) GetUserOrders(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Simular órdenes
	SuccessResponse(c, gin.H{
		"user_id": id,
		"orders": []gin.H{
			{"id": 1, "total": 99.99, "status": "completed"},
			{"id": 2, "total": 49.99, "status": "processing"},
		},
	})
}

// NotifyUser con concurrencia en el handler
func (h *UserHandler) NotifyUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	var req struct {
		Message string `json:"message" binding:"required"`
		Type    string `json:"type" binding:"required,oneof=email sms push"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Procesamiento asíncrono desde el handler
	go func(userID int, msg, msgType string) {
		// Aquí iría el envío real de notificación
		// usando el servicio de notificaciones
		// logger.Printf("Notificación enviada al usuario %d", userID)
	}(id, req.Message, req.Type)

	SuccessResponse(c, gin.H{
		"user_id": id,
		"message": "Notification queued",
		"type":    req.Type,
	})
}
