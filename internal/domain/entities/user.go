package entities

import (
	"errors"
	"fmt"
	"time"
	"user-management/internal/domain/valueobjects"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrWeakPassword = errors.New("password too weak")
)

type User struct {
	ID        uuid.UUID                 `json:"id"`
	Name      string                    `json:"name"`
	Email     string                    `json:"email"`
	Password  valueobjects.PasswordHash `json:"-"`
	Age       int                       `json:"age"`
	Active    bool                      `json:"active"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

func NewUser(name, email string, age int, password string) (usr *User, err error) {
	emailVO, err := valueobjects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := valueobjects.NewPasswordHash(password)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     emailVO.Value(),
		Password:  passwordHash,
		Age:       age,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = t
}

func (u *User) Validate() error {
	if u.Name == "" || len(u.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if u.Age < 0 || u.Age > 120 {
		return fmt.Errorf("age must be between 0 and 120")
	}
	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
}

func (u *User) String() string {
	status := "inactive"
	if u.Active {
		status = "active"
	}
	return fmt.Sprintf("User[ID=%s, Name=%s, Email=%s, Age=%d, Status=%s, CreatedAt=%s, UpdatedAt=%s]",
		u.ID.String(),
		u.Name,
		u.Email,
		u.Age,
		status,
		u.CreatedAt.Format(time.RFC3339),
		u.UpdatedAt.Format(time.RFC3339))
}

func (u *User) Update(name, email string, age int, active bool) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if email != "" {
		emailVO, err := valueobjects.NewEmail(email)
		if err != nil {
			return err
		}
		u.Email = emailVO.Value()
	}

	if age >= 0 && age <= 120 {
		u.Age = age
	}
	u.Active = active

	u.UpdatedAt = time.Now()

	return u.Validate()
}
