package messages

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	domain "server/internal/domain/messages"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func toDomainMessage(row models.ChatMessage) domain.Message {
	return domain.Message{
		ID:             row.ID,
		RoomID:         row.RoomID,
		SenderID:       row.SenderID,
		ReceiverID:     row.ReceiverID,
		OrganisationID: row.OrganisationID,
		Message:        row.Message,
		MessageType:    row.MessageType,
		FileURL:        row.FileURL,
		FileName:       row.FileName,
		FileSize:       row.FileSize,
		IsSeen:         row.IsSeen,
		IsDeleted:      row.IsDeleted,
		ReplyToID:      row.ReplyToID,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
}

func (r *GormRepository) Create(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	m := models.ChatMessage{
		RoomID:         msg.RoomID,
		SenderID:       msg.SenderID,
		ReceiverID:     msg.ReceiverID,
		OrganisationID: msg.OrganisationID,
		Message:        msg.Message,
		MessageType:    msg.MessageType,
		FileURL:        msg.FileURL,
		FileName:       msg.FileName,
		FileSize:       msg.FileSize,
		IsSeen:         msg.IsSeen,
		IsDeleted:      msg.IsDeleted,
		ReplyToID:      msg.ReplyToID,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	res := toDomainMessage(m)
	return &res, nil
}

func (r *GormRepository) ListByRoomCursor(ctx context.Context, roomID string, cursor string, limit int) (*domain.CursorPage, error) {
	if limit <= 0 {
		limit = 30
	}

	var rows []models.ChatMessage
	query := r.db.WithContext(ctx).Where("room_id = ? AND is_deleted = false", roomID)

	if cursor != "" {
		decodedBytes, err := base64.StdEncoding.DecodeString(cursor)
		if err == nil {
			decodedStr := string(decodedBytes)
			parts := strings.Split(decodedStr, ":")
			if len(parts) == 2 {
				id, err1 := strconv.ParseUint(parts[0], 10, 64)
				nanoValue, err2 := strconv.ParseInt(parts[1], 10, 64)
				if err1 == nil && err2 == nil {
					createdAt := time.Unix(0, nanoValue)
					query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)", createdAt, createdAt, id)
				}
			}
		}
	}

	// Fetch limit + 1 to check if there's more
	if err := query.Order("created_at DESC, id DESC").Limit(limit + 1).Find(&rows).Error; err != nil {
		return nil, err
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	// Fetch sender details
	var userIDs []string
	for _, row := range rows {
		userIDs = append(userIDs, row.SenderID)
	}

	userMap := make(map[string]models.UserModel)
	if len(userIDs) > 0 {
		var users []models.UserModel
		r.db.WithContext(ctx).Where("id IN ?", userIDs).Find(&users)
		for _, u := range users {
			userMap[fmt.Sprintf("%d", u.ID)] = u
		}
	}

	var out []domain.Message
	for _, row := range rows {
		msg := toDomainMessage(row)
		if u, ok := userMap[row.SenderID]; ok {
			msg.SenderFirstName = u.FirstName
			msg.SenderLastName = u.LastName
			
			// ProfilePic might be a pointer or string in your struct, copying it
			if u.ProfilePic != "" {
				pic := u.ProfilePic
				msg.SenderAvatarURL = &pic
			}
		}
		out = append(out, msg)
	}

	var nextCursor string
	if len(rows) > 0 {
		lastMsg := rows[len(rows)-1]
		cursorStr := fmt.Sprintf("%d:%d", lastMsg.ID, lastMsg.CreatedAt.UnixNano())
		nextCursor = base64.StdEncoding.EncodeToString([]byte(cursorStr))
	}

	return &domain.CursorPage{
		Messages:   out,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *GormRepository) MarkSeen(ctx context.Context, roomID string, receiverID string) error {
	return r.db.WithContext(ctx).Model(&models.ChatMessage{}).
		Where("room_id = ? AND receiver_id = ? AND is_seen = false", roomID, receiverID).
		Update("is_seen", true).Error
}

func (r *GormRepository) MarkSeenUpTo(ctx context.Context, roomID string, receiverID string, upToMessageID uint) error {
	return r.db.WithContext(ctx).Model(&models.ChatMessage{}).
		Where("room_id = ? AND receiver_id = ? AND id <= ? AND is_seen = false", roomID, receiverID, upToMessageID).
		Update("is_seen", true).Error
}

func (r *GormRepository) DeleteMessage(ctx context.Context, messageID uint, requestingUserID string) error {
	return r.db.WithContext(ctx).Model(&models.ChatMessage{}).
		Where("id = ? AND sender_id = ?", messageID, requestingUserID).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"message":    "[Message deleted]",
		}).Error
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (*domain.Message, error) {
	var m models.ChatMessage
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	res := toDomainMessage(m)
	return &res, nil
}
