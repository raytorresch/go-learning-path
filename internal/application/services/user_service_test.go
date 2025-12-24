package services_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"user-management/internal/application/services"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/ports/input"
	"user-management/internal/domain/valueobjects"
	"user-management/tests/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_RegisterUser(t *testing.T) {
	t.Run("success - creates new user", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		name := "John Doe"
		email := "john@example.com"
		age := 30
		password := "SecurePass123!"

		// Configurar expectativas
		mockRepo.On("FindByEmail", ctx, email).
			Return((*entities.User)(nil), nil).
			Once()

		mockRepo.On("Save", ctx, mock.AnythingOfType("entities.User")).
			Run(func(args mock.Arguments) {
				user := args.Get(1).(entities.User)
				assert.Equal(t, name, user.Name)
				assert.Equal(t, email, user.Email)
				assert.Equal(t, age, user.Age)
				assert.NotEqual(t, uuid.Nil, user.ID)
				assert.True(t, user.Active)
			}).
			Return(nil).
			Once()

		// Act
		user, err := service.RegisterUser(ctx, name, email, age, password)

		// Assert
		assert.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, name, user.Name)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, age, user.Age)
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.True(t, user.Active)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - email already exists", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		existingUser := &entities.User{
			ID:    uuid.New(),
			Name:  "Existing User",
			Email: "existing@example.com",
			Age:   25,
		}

		mockRepo.On("FindByEmail", ctx, "existing@example.com").
			Return(existingUser, nil).
			Once()

		// Act
		user, err := service.RegisterUser(ctx, "New User", "existing@example.com", 30, "Password123!")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, services.ErrEmailAlreadyExists, err)
		assert.Contains(t, err.Error(), "email already exists")

		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Save")
	})

	t.Run("failure - error finding email", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		expectedErr := errors.New("database error")

		mockRepo.On("FindByEmail", ctx, "test@example.com").
			Return((*entities.User)(nil), expectedErr).
			Once()

		log.Println("Mock setup complete for FindByEmail with error")

		// Act
		user, err := service.RegisterUser(ctx, "Test User", "test@example.com", 30, "Password123!")

		log.Println("RegisterUser called, checking results")
		log.Printf("Received error: %v", err)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Save")
	})

	t.Run("failure - invalid user data", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()

		// Email inválido hará que entities.NewUser falle
		mockRepo.On("FindByEmail", ctx, "invalid-email").
			Return((*entities.User)(nil), nil).
			Once()

		// Act
		user, err := service.RegisterUser(ctx, "Test User", "invalid-email", 30, "Password123!")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid email")

		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Save")
	})

	t.Run("failure - save error", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		expectedErr := errors.New("save failed")

		mockRepo.On("FindByEmail", ctx, "test@example.com").
			Return((*entities.User)(nil), nil).
			Once()

		mockRepo.On("Save", ctx, mock.AnythingOfType("entities.User")).
			Return(expectedErr).
			Once()

		// Act
		user, err := service.RegisterUser(ctx, "Test User", "test@example.com", 30, "Password123!")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserProfile(t *testing.T) {
	t.Run("success - returns user profile", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		userID := uuid.New()
		expectedUser := &entities.User{
			ID:    userID,
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   30,
		}

		mockRepo.On("FindByID", ctx, userID).
			Return(expectedUser, nil).
			Once()

		// Act
		user, err := service.GetUserProfile(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - user not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		userID := uuid.New()
		expectedErr := errors.New("user not found")

		mockRepo.On("FindByID", ctx, userID).
			Return((*entities.User)(nil), expectedErr).
			Once()

		// Act
		user, err := service.GetUserProfile(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetAllUsers(t *testing.T) {
	t.Run("success - returns all users", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		expectedUsers := []*entities.User{
			{ID: uuid.New(), Name: "User 1", Email: "user1@example.com"},
			{ID: uuid.New(), Name: "User 2", Email: "user2@example.com"},
		}

		mockRepo.On("GetAllUsers", ctx).
			Return(expectedUsers, nil).
			Once()

		// Act
		users, err := service.GetAllUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Len(t, users, 2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - error getting users", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		expectedErr := errors.New("database error")

		mockRepo.On("GetAllUsers", ctx).
			Return(([]*entities.User)(nil), expectedErr).
			Once()

		// Act
		users, err := service.GetAllUsers(ctx)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("success - returns empty slice when no users", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()

		mockRepo.On("GetAllUsers", ctx).
			Return([]*entities.User{}, nil).
			Once()

		// Act
		users, err := service.GetAllUsers(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Empty(t, users)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	t.Run("success - updates user profile", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		newPasswordHash, err := valueobjects.NewPasswordHash("NewPassword123!")
		if err != nil {
			t.Fatalf("failed to create password hash: %v", err)
		}

		ctx := context.Background()
		user := &entities.User{
			ID:       uuid.New(),
			Name:     "Updated Name",
			Email:    "updated@example.com",
			Age:      35,
			Password: newPasswordHash,
			Active:   true,
		}

		mockRepo.On("Update", ctx, user).
			Return(nil).
			Once()

		// Act
		err = service.UpdateProfile(ctx, user)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - error updating user", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()
		user := &entities.User{
			ID:    uuid.New(),
			Name:  "Test User",
			Email: "test@example.com",
		}
		expectedErr := errors.New("update failed")

		mockRepo.On("Update", ctx, user).
			Return(expectedErr).
			Once()

		// Act
		err := service.UpdateProfile(ctx, user)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("failure - nil user", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)
		service := services.NewUserService(mockRepo)

		ctx := context.Background()

		// Act
		err := service.UpdateProfile(ctx, nil)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user cannot be nil")

		mockRepo.AssertNotCalled(t, "Update")
	})
}

// func TestUserService_DeleteUser(t *testing.T) {
// 	t.Run("success - deletes user", func(t *testing.T) {
// 		// Arrange
// 		mockRepo := new(mocks.MockUserRepository)
// 		service := services.NewUserService(mockRepo)

// 		ctx := context.Background()
// 		userID := uuid.New()

// 		mockRepo.On("Delete", ctx, userID).
// 			Return(nil).
// 			Once()

// 		// Act
// 		err := service.DeleteUser(ctx, userID.String())

// 		// Assert
// 		assert.NoError(t, err)

// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("failure - invalid UUID format", func(t *testing.T) {
// 		// Arrange
// 		mockRepo := new(mocks.MockUserRepository)
// 		service := services.NewUserService(mockRepo)

// 		ctx := context.Background()

// 		// Act
// 		err := service.DeleteUser(ctx, "invalid-uuid")

// 		// Assert
// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "invalid UUID")

// 		mockRepo.AssertNotCalled(t, "Delete")
// 	})

// 	t.Run("failure - error deleting user", func(t *testing.T) {
// 		// Arrange
// 		mockRepo := new(mocks.MockUserRepository)
// 		service := services.NewUserService(mockRepo)

// 		ctx := context.Background()
// 		userID := uuid.New()
// 		expectedErr := errors.New("delete failed")

// 		mockRepo.On("Delete", ctx, userID).
// 			Return(expectedErr).
// 			Once()

// 		// Act
// 		err := service.DeleteUser(ctx, userID.String())

// 		// Assert
// 		assert.Error(t, err)
// 		assert.Equal(t, expectedErr, err)

// 		mockRepo.AssertExpectations(t)
// 	})
// }

func TestUserService_InterfaceImplementation(t *testing.T) {
	t.Run("service implements interface", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)

		// Act
		service := services.NewUserService(mockRepo)

		// Assert
		var _ input.UserService = service
		assert.NotNil(t, service)
	})
}

// Test para cubrir el constructor
func TestNewUserService(t *testing.T) {
	t.Run("creates new service with repository", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockUserRepository)

		// Act
		service := services.NewUserService(mockRepo)

		// Assert
		assert.NotNil(t, service)

		// Verificar que es del tipo correcto
		userService, ok := service.(*services.UserService)
		assert.True(t, ok)
		assert.NotNil(t, userService)
	})
}

// Test para validar la firma del error exportado
func TestErrorConstants(t *testing.T) {
	assert.Equal(t, "email already exists", services.ErrEmailAlreadyExists.Error())
}
