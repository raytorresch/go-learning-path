package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response estándar para todas las APIs
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, status int, err error) {
	c.JSON(status, Response{
		Success: false,
		Error:   err.Error(),
	})
}

func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error:   "validation_error",
		Data:    errors,
	})
}
