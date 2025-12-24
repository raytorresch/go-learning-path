package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Inicio
		start := time.Now()
		path := c.Request.URL.Path

		// Procesar request
		c.Next()

		// Log después de procesar
		latency := time.Since(start)
		status := c.Writer.Status()

		fmt.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
		)
	}
}
