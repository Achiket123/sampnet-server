package messages

import (
	"context"
	domain "server/internal/domain/messages"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) Create(ctx context.Context, msg *domain.Message) error {
	m := models.ChatMessage{RoomID: msg.RoomID, SenderID: msg.SenderID, ReceiverID: msg.ReceiverID, Message: msg.Message, IsSeen: msg.IsSeen, TimeStamp: msg.TimeStamp}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	msg.ID = m.ID
	return nil
}

func (r *GormRepository) ListByRoom(ctx context.Context, roomID string) ([]domain.Message, error) {
	var rows []models.ChatMessage
	if err := r.db.WithContext(ctx).Where("room_id = ?", roomID).Order("time_stamp asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Message, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Message{ID: row.ID, RoomID: row.RoomID, SenderID: row.SenderID, ReceiverID: row.ReceiverID, Message: row.Message, IsSeen: row.IsSeen, TimeStamp: row.TimeStamp})
	}
	return out, nil
}

func (r *GormRepository) MarkSeen(ctx context.Context, roomID string, receiverID string) error {
	return r.db.WithContext(ctx).Model(&models.ChatMessage{}).Where("room_id = ? AND receiver_id = ?", roomID, receiverID).Update("is_seen", true).Error
}
