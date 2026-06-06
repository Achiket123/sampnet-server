package notifications

import (
	"context"
	"server/internal/domain/notifications"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) notifications.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, n *notifications.Notification) error {
	model := toModel(n)
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now()
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *gormRepository) GetByUser(ctx context.Context, userID uint, offset int) ([]notifications.Notification, error) {
	var ms []models.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(20).
		Offset(offset * 20).
		Find(&ms).Error

	if err != nil {
		return nil, err
	}

	var ns []notifications.Notification
	for _, m := range ms {
		ns = append(ns, toDomain(&m))
	}
	return ns, nil
}

func (r *gormRepository) MarkRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *gormRepository) MarkAllRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func toModel(n *notifications.Notification) *models.Notification {
	return &models.Notification{
		Model: gorm.Model{
			ID:        n.ID,
			CreatedAt: n.CreatedAt,
		},
		UserID:           n.UserID,
		OrganisationID:   n.OrganisationID,
		Title:            n.Title,
		Message:          n.Message,
		IsRead:           n.IsRead,
		NotificationType: n.NotificationType,
		Link:             n.Link,
	}
}

func toDomain(m *models.Notification) notifications.Notification {
	return notifications.Notification{
		ID:               m.ID,
		UserID:           m.UserID,
		OrganisationID:   m.OrganisationID,
		Title:            m.Title,
		Message:          m.Message,
		IsRead:           m.IsRead,
		NotificationType: m.NotificationType,
		Link:             m.Link,
		CreatedAt:        m.CreatedAt,
	}
}
