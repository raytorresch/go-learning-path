package services

import (
	"errors"
	"testing"
	"user-management/internal/domain/models"
)

// Test doubles
type mockUserRepo struct {
	saveFunc func(*models.User) (*models.User, error)
	findFunc func(int) (*models.User, error)
}

func (m *mockUserRepo) Save(user *models.User) (*models.User, error) {
	return m.saveFunc(user)
}

func (m *mockUserRepo) FindByID(id int) (*models.User, error) {
	return m.findFunc(id)
}

func (m *mockUserRepo) FindAll() ([]*models.User, error) {
	return nil, nil
}

func (m *mockUserRepo) Update(user *models.User) (*models.User, error) {
	return nil, nil
}

func (m *mockUserRepo) Delete(id int) error {
	return nil
}

func (m *mockUserRepo) GetUserByID(id int) (int, error) {
	return 0, nil
}

// Implementa otros métodos...

func TestUserService_RegisterUser(t *testing.T) {
	tests := []struct {
		name       string
		setupMock  func() *mockUserRepo
		wantErr    bool
		wantUserID int
	}{
		{
			name: "registro exitoso",
			setupMock: func() *mockUserRepo {
				return &mockUserRepo{
					saveFunc: func(user *models.User) (*models.User, error) {
						user.ID = 1 // Simular ID asignado
						return user, nil
					},
				}
			},
			wantErr:    false,
			wantUserID: 1,
		},
		{
			name: "error en repository",
			setupMock: func() *mockUserRepo {
				return &mockUserRepo{
					saveFunc: func(user *models.User) (*models.User, error) {
						return nil, errors.New("database error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			service := NewUserService(mockRepo, nil)

			user, err := service.RegisterUser("Test", "test@test.com", 30)

			if tt.wantErr {
				if err == nil {
					t.Error("Esperaba error, obtuve nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("RegisterUser falló: %v", err)
			}

			if user.ID != int64(tt.wantUserID) {
				t.Errorf("RegisterUser() ID = %d, quiere %d", user.ID, tt.wantUserID)
			}
		})
	}
}

// func TestUserService_GetUserStats(t *testing.T) {
//     mockRepo := &mockUserRepo{
//         findFunc: func(id int) (*models.User, error) {
//             return &models.User{ID: int64(id), Active: true}, nil
//         },
//     }

//     service := NewUserService(mockRepo, nil)

//     // Test básico
//     user, err := service.GetUserByID(1)
//     if err != nil {
//         t.Fatalf("GetUserByID falló: %v", err)
//     }

//     if user.ID != 1 {
//         t.Errorf("GetUserByID() = %d, quiere 1", user.ID)
//     }
// }
