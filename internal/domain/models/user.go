package models

import (
	"fmt"
	"time"
	"user-management/pkg/utils"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	email     string    `json:"email"`
	Age       int       `json:"age"`
	Active    bool      `json:"active"`
	createdAt time.Time `json:"created_at"`
	updatedAt time.Time `json:"updated_at"`
}

func NewUser(name, email string, age int, active bool) *User {
	return &User{
		Name:      name,
		email:     email,
		Age:       age,
		Active:    active,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func (u *User) Email() string {
	return u.email
}

func (u *User) SetEmail(email string) error {
	if !utils.ValidateEmail(email) {
		return fmt.Errorf("email inválido")
	}
	return nil
}

func (u *User) SetUpdatedAt(t time.Time) {
	u.updatedAt = t
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if u.email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if u.Age < 0 || u.Age > 120 {
		return fmt.Errorf("age must be between 0 and 120")
	}

	return nil
}

func (u *User) String() string {
	status := "inactive"
	if u.Active {
		status = "active"
	}
	return fmt.Sprintf("User[ID=%d, Name=%s, Email=%s, Age=%d, Status=%s, CreatedAt=%s, UpdatedAt=%s]",
		u.ID, u.Name, u.email, u.Age, status, u.createdAt.Format(time.RFC3339), u.updatedAt.Format(time.RFC3339))
}

func (u *User) Update(name, email string, age int, active bool) error {
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.email = email
	}
	if age >= 0 && age <= 120 {
		u.Age = age
	}
	u.Active = active

	u.updatedAt = time.Now()

	return u.Validate()
}
