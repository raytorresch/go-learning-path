package entities_test

import (
	"testing"
	"time"
	"user-management/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestNewUser() {
	tests := []struct {
		name        string
		email       string
		age         int
		userName    string
		password    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "success - valid user",
			email:    "test@example.com",
			userName: "John Doe",
			age:      30,
			password: "SecurePass123!",
			wantErr:  false,
		},
		{
			name:        "fail - invalid email",
			email:       "not-an-email",
			age:         30,
			userName:    "John",
			password:    "password123",
			wantErr:     true,
			errContains: "invalid email format",
		},
		{
			name:        "fail - short username",
			email:       "test@mail.com",
			age:         150,
			userName:    "J", // Muy corto
			password:    "password123",
			wantErr:     true,
			errContains: "name must be at least 2 characters",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			user, err := entities.NewUser(tt.userName, tt.email, tt.age, tt.password)

			if tt.wantErr {
				s.Error(err)
				if tt.errContains != "" {
					s.Contains(err.Error(), tt.errContains)
				}
				s.Nil(user)
			} else {
				s.NoError(err)
				s.NotNil(user)
				s.Equal(tt.email, user.Email)
				s.Equal(tt.userName, user.Name)
				s.Equal(tt.age, user.Age)
				s.False(user.CreatedAt.IsZero())
				s.NotEqual(uuid.Nil, user.ID)
			}
		})
	}
}

func (s *UserTestSuite) TestUser_UpdateProfile() {
	user, err := entities.NewUser("John Updated", "test@example.com", 43, "Password123!")
	s.NoError(err)
	s.NotNil(user)

	// Test success
	err = user.Update("John Updated", "example@example.com", 20, false)
	s.NoError(err)
	s.Equal("John Updated", user.Name)
	s.Equal("example@example.com", user.Email)
	s.Equal(20, user.Age)
	s.False(user.Active)

	// Test failure
	err = user.Update("", "example@example.com", 121, false) // Nombre vacío
	s.Error(err)
	s.Contains(err.Error(), "name cannot be empty")
}

func (s *UserTestSuite) TestUser_Validate() {
	tests := []struct {
		name        string
		user        *entities.User
		wantErr     bool
		errContains string
	}{
		{
			name: "success - valid user",
			user: &entities.User{
				ID:        uuid.New(),
				Name:      "Valid User",
				Email:     "valid@example.com",
				Age:       30,
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "fail - empty name",
			user: &entities.User{
				ID:    uuid.New(),
				Name:  "", // Nombre vacío
				Email: "valid@example.com",
				Age:   30,
			},
			wantErr:     true,
			errContains: "name must be at least 2 characters",
		},
		{
			name: "fail - invalid age",
			user: &entities.User{
				ID:    uuid.New(),
				Name:  "User",
				Email: "valid@example.com",
				Age:   150, // Edad inválida
			},
			wantErr:     true,
			errContains: "age must be between 0 and 120",
		},
		{
			name: "fail - empty email",
			user: &entities.User{
				ID:    uuid.New(),
				Name:  "User",
				Email: "", // Email vacío
				Age:   30,
			},
			wantErr:     true,
			errContains: "email cannot be empty",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := tt.user.Validate()
			if tt.wantErr {
				s.Error(err)
				if tt.errContains != "" {
					s.Contains(err.Error(), tt.errContains)
				}
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *UserTestSuite) TestUser_String() {
	user, err := entities.NewUser("Alice", "valid@example.com", 30, "Password123!")
	s.NoError(err)
	s.NotNil(user)

	userStr := user.String()
	s.Contains(userStr, "Alice")
	s.Contains(userStr, "valid@example.com")
	s.Contains(userStr, "30")
	s.Contains(userStr, "active")
}

func (s *UserTestSuite) TestUser_SetUpdatedAt() {
	user, err := entities.NewUser("Bob", "valid@example.com", 30, "Password123!")
	s.NoError(err)
	s.NotNil(user)

	newTime := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	user.SetUpdatedAt(newTime)

	s.Equal(newTime, user.UpdatedAt)
}
