package notifications

import (
	"context"
)

type UseCase interface {
	CreateNotification(ctx context.Context, userID uint, organisationID uint, title string, message string, notificationType string, link string) error
	GetNotifications(ctx context.Context, userID uint, offset int) ([]Notification, error)
	MarkNotificationRead(ctx context.Context, notificationID uint) error
	MarkAllNotificationsRead(ctx context.Context, userID uint) error
}
