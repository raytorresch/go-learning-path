package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"user-management/internal/domain/entities"
	"user-management/internal/domain/ports/input"
	"user-management/internal/domain/valueobjects"
)

type OrderHandler struct {
	orderService input.OrderService
}

func NewOrderHandler(orderService input.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/orders", h.CreateOrder)
	router.GET("/orders/:id", h.GetOrder)
	router.GET("/orders", h.ListOrders)
	router.POST("/orders/:id/cancel", h.CancelOrder)
	router.GET("/orders/:id/stream", h.StreamOrderEvents) // Server-Sent Events
}

// CreateOrder con validación compleja
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order entities.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Simular creación
	order.ID = 1
	order.Status = valueobjects.OrderStatus(entities.StatusProcessing)

	// Procesar concurrentemente
	go func(o entities.Order) {
		// Procesamiento en background
		// h.orderService.ProcessOrder(&o)
	}(order)

	c.JSON(http.StatusAccepted, Response{
		Success: true,
		Data:    order,
		Message: "Order processing started",
	})
}

// GetOrder - Obtener pedido por ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Mock búsqueda
	order := &entities.Order{
		ID:     id,
		Total:  99.99,
		Status: valueobjects.OrderStatus(entities.StatusCompleted),
	}

	SuccessResponse(c, order)
}

// ListOrders - Listar pedidos con paginación
func (h *OrderHandler) ListOrders(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	// Mock listado
	orders := []entities.Order{}
	for i := range limit {
		orderID := (page-1)*limit + i + 1
		orders = append(orders, entities.Order{
			ID:     orderID,
			Total:  float64(20 * orderID),
			Status: valueobjects.OrderStatus(entities.StatusCompleted),
		})
	}

	SuccessResponse(c, gin.H{
		"page":   page,
		"limit":  limit,
		"orders": orders,
	})
}

// CancelOrder - Cancelar un pedido
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// Simular cancelación
	order := &entities.Order{
		ID:     id,
		Status: valueobjects.OrderStatus(entities.StatusCancelled),
	}

	SuccessResponse(c, gin.H{
		"message": fmt.Sprintf("Order %d cancelled", id),
		"order":   order,
	})
}

// StreamOrderEvents - Server-Sent Events
func (h *OrderHandler) StreamOrderEvents(c *gin.Context) {
	idStr := c.Param("id")
	orderID, _ := strconv.Atoi(idStr)

	// Configurar SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Channel para eventos
	eventChan := make(chan string)

	// Simular eventos periódicos
	go func() {
		defer close(eventChan)

		events := []string{"processing", "packing", "shipping", "delivered"}
		for _, event := range events {
			select {
			case <-c.Request.Context().Done():
				return
			case eventChan <- event:
				time.Sleep(2 * time.Second)
			}
		}
	}()

	// Stream eventos al cliente
	c.Stream(func(w io.Writer) bool {
		if event, ok := <-eventChan; ok {
			c.SSEvent("message", gin.H{
				"order_id": orderID,
				"event":    event,
				"time":     time.Now().Unix(),
			})
			return true
		}
		return false
	})
}
