package services

import "user-management/internal/domain/models"

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
