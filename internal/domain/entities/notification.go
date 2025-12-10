package entities

import "time"

// Notification representa una notificación base
type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Sent      bool      `json:"sent"`
	CreatedAt time.Time `json:"created_at"`
}

// Notifier interface define el contrato para enviar notificaciones
type Notifier interface {
	Send(notification *Notification) error
	GetType() string
}

// EmailNotification implementa Notifier
type EmailNotification struct {
	Notification
	Email string `json:"email"`
}

func (e *EmailNotification) Send(notification *Notification) error {
	// Simular envío de email
	notification.Sent = true
	return nil
}

func (e *EmailNotification) GetType() string {
	return "email"
}

// SMSNotification implementa Notifier
type SMSNotification struct {
	Notification
	Phone string `json:"phone"`
}

func (s *SMSNotification) Send(notification *Notification) error {
	// Simular envío de SMS
	notification.Sent = true
	return nil
}

func (s *SMSNotification) GetType() string {
	return "sms"
}
