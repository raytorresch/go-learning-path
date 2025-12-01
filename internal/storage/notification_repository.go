package storage

import (
	"sync"
	"time"
	"user-management/internal/models"
)

type NotificationRepository struct {
	notifications map[int]*models.Notification
	mutex         sync.RWMutex
	nextID        int
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		notifications: make(map[int]*models.Notification),
		nextID:        1,
	}
}

func (r *NotificationRepository) Save(notification *models.Notification) (*models.Notification, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	notification.ID = r.nextID
	notification.CreatedAt = time.Now()
	r.notifications[notification.ID] = notification
	r.nextID++

	return notification, nil
}

func (r *NotificationRepository) FindByUserID(userID int) []*models.Notification {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var userNotifications []*models.Notification
	for _, notification := range r.notifications {
		if notification.UserID == userID {
			userNotifications = append(userNotifications, notification)
		}
	}
	return userNotifications
}
