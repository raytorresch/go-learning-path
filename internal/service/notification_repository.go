package service

import "user-management/internal/models"

type NotificationService struct {
	notifiers []models.Notifier
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifiers: []models.Notifier{
			&models.EmailNotification{},
			&models.SMSNotification{},
		},
	}
}

// SendToAllChannels demuestra polimorfismo
func (s *NotificationService) SendToAllChannels(notification *models.Notification, user *models.User) error {
	for _, notifier := range s.notifiers {
		// Type assertion para obtener implementación específica
		switch n := notifier.(type) {
		case *models.EmailNotification:
			n.Email = user.Email
			n.Send(notification)
		case *models.SMSNotification:
			n.Phone = "+1234567890" // Simulado
			n.Send(notification)
		}
	}
	return nil
}

// GetAvailableNotifiers muestra interfaces en acción
func (s *NotificationService) GetAvailableNotifiers() []string {
	var types []string
	for _, notifier := range s.notifiers {
		types = append(types, notifier.GetType())
	}
	return types
}
