package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-management/internal/application/services"
	"user-management/internal/domain/entities"
	"user-management/internal/infrastructure/http/handlers"
	"user-management/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Active    bool      `json:"active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - creates user", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)

		mockRepo.On("FindByEmail", mock.Anything, "test@example.com").
			Return((*entities.User)(nil), nil)

		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("entities.User")).
			Run(func(args mock.Arguments) {
				user := args.Get(1).(entities.User)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "John Doe", user.Name)
			}).
			Return(nil)

		userService := services.NewUserService(mockRepo)
		handler := handlers.NewUserHandler(userService)

		requestBody := map[string]any{
			"name":     "John Doe",
			"email":    "test@example.com",
			"age":      30,
			"password": "SecurePass123!",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		// ctx, engine := gin.CreateTestContext(w)

		router := gin.New()
		apiGroup := router.Group("/")
		handler.RegisterRoutes(apiGroup)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var apiResp ApiResponse
		err := json.Unmarshal(w.Body.Bytes(), &apiResp)
		assert.NoError(t, err)
		assert.True(t, apiResp.Success)

		// Convertir el data a UserResponse
		dataBytes, _ := json.Marshal(apiResp.Data)
		var userResp UserResponse
		err = json.Unmarshal(dataBytes, &userResp)
		assert.NoError(t, err)

		assert.Equal(t, "test@example.com", userResp.Email)
		assert.Equal(t, "John Doe", userResp.Name)
		assert.Equal(t, 30, userResp.Age)
		assert.True(t, userResp.Active)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - email already exists", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)

		// Simular que el email ya existe
		existingUser := &entities.User{
			ID:    uuid.New(),
			Email: "test@example.com",
			Name:  "Existing User",
			Age:   40,
		}

		mockRepo.On("FindByEmail", mock.Anything, "test@example.com").
			Return(existingUser, nil)

		userService := services.NewUserService(mockRepo)
		handler := handlers.NewUserHandler(userService)

		requestBody := map[string]any{
			"email":    "test@example.com",
			"name":     "John Doe",
			"age":      30,
			"password": "SecurePass123!",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		// ctx, _ := gin.CreateTestContext(w)

		router := gin.New()
		apiGroup := router.Group("/")
		handler.RegisterRoutes(apiGroup)

		router.ServeHTTP(w, req)

		t.Logf("Status Code: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())
		assert.Equal(t, http.StatusConflict, w.Code)

		var apiResp ApiResponse
		err := json.Unmarshal(w.Body.Bytes(), &apiResp)
		assert.NoError(t, err)
		assert.False(t, apiResp.Success)
		assert.Contains(t, apiResp.Error, "already")

		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Create")
	})
}
