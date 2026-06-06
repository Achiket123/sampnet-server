package notifications

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByUser(ctx context.Context, userID uint, offset int) ([]Notification, error)
	MarkRead(ctx context.Context, notificationID uint) error
	MarkAllRead(ctx context.Context, userID uint) error
}
