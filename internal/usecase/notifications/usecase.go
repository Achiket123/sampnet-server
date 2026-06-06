package notifications

import (
	"context"
	"server/internal/domain/notifications"
	"time"
)

type service struct {
	repo        notifications.Repository
	broadcaster notifications.Broadcaster
}

func NewService(repo notifications.Repository, broadcaster notifications.Broadcaster) notifications.UseCase {
	return &service{
		repo:        repo,
		broadcaster: broadcaster,
	}
}

func (s *service) CreateNotification(ctx context.Context, userID uint, organisationID uint, title string, message string, notificationType string, link string) error {
	n := &notifications.Notification{
		UserID:           userID,
		OrganisationID:   organisationID,
		Title:            title,
		Message:          message,
		IsRead:           false,
		NotificationType: notificationType,
		Link:             link,
		CreatedAt:        time.Now(),
	}
	err := s.repo.Create(ctx, n)
	if err != nil {
		return err
	}
	if s.broadcaster != nil {
		_ = s.broadcaster.BroadcastNotification(ctx, n)
	}
	return nil
}

func (s *service) GetNotifications(ctx context.Context, userID uint, offset int) ([]notifications.Notification, error) {
	return s.repo.GetByUser(ctx, userID, offset)
}

func (s *service) MarkNotificationRead(ctx context.Context, notificationID uint) error {
	return s.repo.MarkRead(ctx, notificationID)
}

func (s *service) MarkAllNotificationsRead(ctx context.Context, userID uint) error {
	return s.repo.MarkAllRead(ctx, userID)
}

