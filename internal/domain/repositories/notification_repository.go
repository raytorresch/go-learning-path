package repositories

import (
	"sync"
	"time"
	"user-management/internal/domain/entities"
)

type NotificationRepository struct {
	notifications map[int]*entities.Notification
	mutex         sync.RWMutex
	nextID        int
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		notifications: make(map[int]*entities.Notification),
		nextID:        1,
	}
}

func (r *NotificationRepository) Save(notification *entities.Notification) (*entities.Notification, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	notification.ID = r.nextID
	notification.CreatedAt = time.Now()
	r.notifications[notification.ID] = notification
	r.nextID++

	return notification, nil
}

func (r *NotificationRepository) FindByUserID(userID int) []*entities.Notification {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var userNotifications []*entities.Notification
	for _, notification := range r.notifications {
		if notification.UserID == userID {
			userNotifications = append(userNotifications, notification)
		}
	}
	return userNotifications
}
