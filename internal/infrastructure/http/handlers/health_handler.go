package handlers

import (
	"runtime"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/health", h.HealthCheck)
	router.GET("/metrics", h.SystemMetrics)
}

// HealthCheck muestra estado del sistema
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	SuccessResponse(c, gin.H{
		"status":  "healthy",
		"service": "order-management",
		"version": "1.0.0",
	})
}

// SystemMetrics muestra métricas básicas
func (h *HealthHandler) SystemMetrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	SuccessResponse(c, gin.H{
		"memory": gin.H{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"num_gc":      m.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
	})
}
