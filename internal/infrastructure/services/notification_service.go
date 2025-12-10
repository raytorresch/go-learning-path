package services

import "user-management/internal/domain/entities"

type NotificationService struct {
	notifiers []entities.Notifier
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifiers: []entities.Notifier{
			&entities.EmailNotification{},
			&entities.SMSNotification{},
		},
	}
}
