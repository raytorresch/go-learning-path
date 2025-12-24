package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
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

func TestUserHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - get user by ID", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)

		userID := uuid.New()
		existingUser := &entities.User{
			ID:     userID,
			Email:  "test@example.com",
			Name:   "Existing User",
			Age:    40,
			Active: true,
		}

		mockRepo.On("FindByID", mock.Anything, userID).
			Return(existingUser, nil)

		userService := services.NewUserService(mockRepo)
		handler := handlers.NewUserHandler(userService)

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)

		w := httptest.NewRecorder()

		router := gin.New()
		apiGroup := router.Group("/")
		handler.RegisterRoutes(apiGroup)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var apiResp ApiResponse
		err := json.Unmarshal(w.Body.Bytes(), &apiResp)
		assert.NoError(t, err)
		assert.True(t, apiResp.Success)

		dataBytes, _ := json.Marshal(apiResp.Data)
		var userResp UserResponse
		err = json.Unmarshal(dataBytes, &userResp)
		assert.NoError(t, err)

		assert.Equal(t, existingUser.ID, userResp.ID)
		assert.Equal(t, existingUser.Email, userResp.Email)
		assert.Equal(t, existingUser.Name, userResp.Name)
		assert.Equal(t, existingUser.Age, userResp.Age)
		assert.Equal(t, existingUser.Active, userResp.Active)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - user not found", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)

		userID := uuid.New()

		mockRepo.On("FindByID", mock.Anything, userID).
			Return((*entities.User)(nil), errors.New("user not found"))

		userService := services.NewUserService(mockRepo)
		handler := handlers.NewUserHandler(userService)

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)

		w := httptest.NewRecorder()

		router := gin.New()
		apiGroup := router.Group("/")
		handler.RegisterRoutes(apiGroup)

		router.ServeHTTP(w, req)
		t.Logf("Status Code: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())
		assert.Equal(t, http.StatusNotFound, w.Code)

		var apiResp ApiResponse
		err := json.Unmarshal(w.Body.Bytes(), &apiResp)
		assert.NoError(t, err)
		assert.False(t, apiResp.Success)
		assert.Contains(t, apiResp.Error, "not found")

		mockRepo.AssertExpectations(t)
	})
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - get all users", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)

		users := []*entities.User{
			{
				ID:     uuid.New(),
				Email:  "test@example.com",
				Name:   "Existing User",
				Age:    40,
				Active: true,
			},
			{
				ID:     uuid.New(),
				Email:  "test2@example.com",
				Name:   "Second User",
				Age:    25,
				Active: false,
			},
		}

		mockRepo.On("GetAllUsers", mock.Anything).
			Return(users, nil)

		userService := services.NewUserService(mockRepo)
		handler := handlers.NewUserHandler(userService)

		req, _ := http.NewRequest("GET", "/users", nil)

		w := httptest.NewRecorder()

		router := gin.New()
		apiGroup := router.Group("/")
		handler.RegisterRoutes(apiGroup)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var apiResp ApiResponse
		err := json.Unmarshal(w.Body.Bytes(), &apiResp)
		assert.NoError(t, err)
		assert.True(t, apiResp.Success)

		dataBytes, _ := json.Marshal(apiResp.Data)
		var usersResp []UserResponse
		err = json.Unmarshal(dataBytes, &usersResp)
		assert.NoError(t, err)

		assert.Len(t, usersResp, 2)
		assert.Equal(t, users[0].Email, usersResp[0].Email)
		assert.Equal(t, users[1].Email, usersResp[1].Email)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success - update user", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)

		userID := uuid.New()
		existingUser := &entities.User{
			ID:     userID,
			Email:  "test@example.com",
			Name:   "Existing User",
			Age:    40,
			Active: true,
		}

		mockRepo.On("FindByID", mock.Anything, userID).
			Return(existingUser, nil)

		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.User")).
			Run(func(args mock.Arguments) {
				updatedUser := args.Get(1).(*entities.User)
				assert.Equal(t, "Updated User", updatedUser.Name)
				assert.Equal(t, 45, updatedUser.Age)
			}).
			Return(nil)

		userService := services.NewUserService(mockRepo)
		handler := handlers.NewUserHandler(userService)

		requestBody := map[string]any{
			"name":   "Updated User",
			"age":    45,
			"email":  "updated@example.com",
			"active": false,
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest("PUT", "/users/"+existingUser.ID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router := gin.New()
		apiGroup := router.Group("/")
		handler.RegisterRoutes(apiGroup)

		router.ServeHTTP(w, req)
		log.Printf("Status Code: %d", w.Code)
		log.Printf("Response Body: %s", w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)

		var apiResp ApiResponse
		err := json.Unmarshal(w.Body.Bytes(), &apiResp)
		assert.NoError(t, err)
		assert.True(t, apiResp.Success)

		dataBytes, _ := json.Marshal(apiResp.Data)
		var userResp UserResponse
		err = json.Unmarshal(dataBytes, &userResp)
		assert.NoError(t, err)

		assert.Equal(t, existingUser.ID, userResp.ID)
		assert.Equal(t, "Updated User", userResp.Name)
		assert.Equal(t, 45, userResp.Age)
		assert.Equal(t, "updated@example.com", userResp.Email)
		assert.Equal(t, false, userResp.Active)
		mockRepo.AssertExpectations(t)
	})
}
