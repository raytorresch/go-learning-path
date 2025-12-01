package models

import (
	"fmt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if u.Email == "" {
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
		u.ID, u.Name, u.Email, u.Age, status, u.CreatedAt.Format(time.RFC3339), u.UpdatedAt.Format(time.RFC3339))
}

func (u *User) Update(name, email string, age int, active bool) error {
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	if age >= 0 && age <= 120 {
		u.Age = age
	}
	u.Active = active

	u.UpdatedAt = time.Now()

	return u.Validate()
}
