package chats

import (
	"context"
	"fmt"
	"strings"
	domain "server/internal/domain/chats"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func toDomainChat(row models.Chat, participants []models.ChatParticipant) domain.Chat {
	var parts []domain.ChatParticipant
	for _, p := range participants {
		parts = append(parts, domain.ChatParticipant{
			ChatID:            p.ChatID,
			UserID:            p.UserID,
			UnreadCount:       p.UnreadCount,
			LastReadMessageID: p.LastReadMessageID,
			JoinedAt:          p.JoinedAt,
		})
	}
	var createdBy uint
	if row.CreatedBy != nil {
		createdBy = *row.CreatedBy
	}

	return domain.Chat{
		ID:             row.ID,
		RoomID:         row.RoomID,
		OrganisationID: row.OrganisationID,
		Name:           row.Name,
		IsGroup:        row.IsGroup,
		CreatedBy:      createdBy,
		LastMessage:    row.LastMessage,
		LastMessageAt:  row.LastMessageAt,
		MessageCount:   row.MessageCount,
		Participants:   parts,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
}

func (r *GormRepository) Create(ctx context.Context, chat *domain.Chat) error {
	var createdBy *uint
	if chat.CreatedBy > 0 {
		createdBy = &chat.CreatedBy
	}
	m := models.Chat{
		RoomID:         chat.RoomID,
		OrganisationID: chat.OrganisationID,
		Name:           chat.Name,
		IsGroup:        chat.IsGroup,
		CreatedBy:      createdBy,
		LastMessage:    chat.LastMessage,
		LastMessageAt:  chat.LastMessageAt,
		MessageCount:   chat.MessageCount,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	chat.ID = m.ID
	chat.CreatedAt = m.CreatedAt
	chat.UpdatedAt = m.UpdatedAt

	// If it's a group and RoomID wasn't set, set it now
	if chat.IsGroup && chat.RoomID == "" {
		chat.RoomID = fmt.Sprintf("group_%d", chat.ID)
		r.db.WithContext(ctx).Model(&m).Update("room_id", chat.RoomID)
	}

	// Insert participants
	for i, p := range chat.Participants {
		pModel := models.ChatParticipant{
			ChatID:      m.ID,
			UserID:      p.UserID,
			UnreadCount: p.UnreadCount,
		}
		if err := r.db.WithContext(ctx).Create(&pModel).Error; err != nil {
			return err
		}
		chat.Participants[i].ChatID = m.ID
		chat.Participants[i].JoinedAt = pModel.JoinedAt
	}

	return nil
}

func (r *GormRepository) ListByOrganisation(ctx context.Context, organisationID uint) ([]domain.Chat, error) {
	var rows []models.Chat
	if err := r.db.WithContext(ctx).Where("organisation_id = ?", organisationID).Order("last_message_at desc nulls last").Find(&rows).Error; err != nil {
		return nil, err
	}
	
	// Preload participants manually (or could use gorm preload if we define relations, but we don't have explicit relations in models currently)
	var chatIDs []uint
	for _, row := range rows {
		chatIDs = append(chatIDs, row.ID)
	}

	var allParticipants []models.ChatParticipant
	if len(chatIDs) > 0 {
		r.db.WithContext(ctx).Where("chat_id IN ?", chatIDs).Find(&allParticipants)
	}

	partMap := make(map[uint][]models.ChatParticipant)
	for _, p := range allParticipants {
		partMap[p.ChatID] = append(partMap[p.ChatID], p)
	}

	var out []domain.Chat
	for _, row := range rows {
		out = append(out, toDomainChat(row, partMap[row.ID]))
	}
	return out, nil
}

func (r *GormRepository) ListByUser(ctx context.Context, userID uint, orgID uint) ([]domain.Chat, error) {
	// Auto-create DMs with all employees in the org
	var employees []models.Employee
	if err := r.db.WithContext(ctx).Where("organisation_id = ?", orgID).Find(&employees).Error; err == nil {
		for _, emp := range employees {
			if emp.UserID != userID && emp.UserID != 0 {
				_, _ = r.GetOrCreateDM(ctx, userID, emp.UserID, orgID)
			}
		}
	}

	var participants []models.ChatParticipant
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&participants).Error; err != nil {
		return nil, err
	}

	var chatIDs []uint
	for _, p := range participants {
		chatIDs = append(chatIDs, p.ChatID)
	}

	if len(chatIDs) == 0 {
		return []domain.Chat{}, nil
	}

	var rows []models.Chat
	if err := r.db.WithContext(ctx).Where("id IN ? AND organisation_id = ?", chatIDs, orgID).Order("last_message_at desc nulls last").Find(&rows).Error; err != nil {
		return nil, err
	}

	// fetch all participants for these chats
	var allParticipants []models.ChatParticipant
	if err := r.db.WithContext(ctx).Where("chat_id IN ?", chatIDs).Find(&allParticipants).Error; err == nil {
		// Also fetch user details to populate FirstName, LastName, AvatarURL
		var userIDs []uint
		for _, p := range allParticipants {
			userIDs = append(userIDs, p.UserID)
		}

		var users []models.UserModel
		if len(userIDs) > 0 {
			r.db.WithContext(ctx).Where("id IN ?", userIDs).Find(&users)
		}
		userMap := make(map[uint]models.UserModel)
		for _, u := range users {
			userMap[u.ID] = u
		}

		partMap := make(map[uint][]domain.ChatParticipant)
		for _, p := range allParticipants {
			dp := domain.ChatParticipant{
				ChatID:            p.ChatID,
				UserID:            p.UserID,
				UnreadCount:       p.UnreadCount,
				LastReadMessageID: p.LastReadMessageID,
				JoinedAt:          p.JoinedAt,
			}
			if u, ok := userMap[p.UserID]; ok {
				dp.FirstName = u.FirstName
				dp.LastName = u.LastName
				dp.AvatarURL = u.ProfilePic
			}
			partMap[p.ChatID] = append(partMap[p.ChatID], dp)
		}

		var out []domain.Chat
		for _, row := range rows {
			chat := toDomainChat(row, nil)
			chat.Participants = partMap[row.ID]
			
			// Auto name DM
			if !chat.IsGroup && len(chat.Participants) == 2 {
				for _, p := range chat.Participants {
					if p.UserID != userID {
						firstName := strings.TrimSpace(p.FirstName)
						lastName := strings.TrimSpace(p.LastName)
						if firstName == "" && lastName == "" {
							chat.Name = fmt.Sprintf("User %d", p.UserID)
						} else if firstName == "" {
							chat.Name = lastName
						} else if lastName == "" {
							chat.Name = firstName
						} else {
							chat.Name = fmt.Sprintf("%s %s", firstName, lastName)
						}
						break
					}
				}
			}
			
			out = append(out, chat)
		}
		return out, nil
	}

	return nil, nil
}

func (r *GormRepository) GetByID(ctx context.Context, chatID uint) (*domain.Chat, error) {
	var row models.Chat
	if err := r.db.WithContext(ctx).First(&row, chatID).Error; err != nil {
		return nil, err
	}

	var parts []models.ChatParticipant
	r.db.WithContext(ctx).Where("chat_id = ?", chatID).Find(&parts)

	chat := toDomainChat(row, parts)
	return &chat, nil
}

func (r *GormRepository) GetByRoomID(ctx context.Context, roomID string) (*domain.Chat, error) {
	var row models.Chat
	if err := r.db.WithContext(ctx).Where("room_id = ?", roomID).First(&row).Error; err != nil {
		return nil, err
	}

	var parts []models.ChatParticipant
	r.db.WithContext(ctx).Where("chat_id = ?", row.ID).Find(&parts)

	chat := toDomainChat(row, parts)
	return &chat, nil
}

func (r *GormRepository) UpdateLastMessage(ctx context.Context, id uint, message string, at *string) error {
	var ts *time.Time
	if at != nil && *at != "" {
		parsed, err := time.Parse(time.RFC3339, *at)
		if err == nil {
			ts = &parsed
		}
	}
	updates := map[string]any{"last_message": message, "last_message_at": ts}
	return r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", id).Updates(updates).Error
}

func (r *GormRepository) IncrementMessageCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", id).UpdateColumn("message_count", gorm.Expr("message_count + 1")).Error
}

func (r *GormRepository) GetOrCreateDM(ctx context.Context, userA uint, userB uint, orgID uint) (*domain.Chat, error) {
	// Find existing DM between userA and userB
	var chatID uint
	
	// Query to find a chat where both participants exist AND IsGroup = false
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.id FROM chats c
		JOIN chat_participants p1 ON c.id = p1.chat_id
		JOIN chat_participants p2 ON c.id = p2.chat_id
		WHERE c.is_group = false AND c.organisation_id = ?
		AND p1.user_id = ? AND p2.user_id = ?
	`, orgID, userA, userB).Scan(&chatID).Error

	if err != nil || chatID == 0 {
		var smaller, larger uint
		if userA < userB {
			smaller = userA
			larger = userB
		} else {
			smaller = userB
			larger = userA
		}
		roomID := fmt.Sprintf("dm_%d_%d", smaller, larger)

		// create new DM
		chat := &domain.Chat{
			RoomID:         roomID,
			OrganisationID: orgID,
			IsGroup:        false,
			Participants: []domain.ChatParticipant{
				{UserID: userA},
				{UserID: userB},
			},
		}
		if err := r.Create(ctx, chat); err != nil {
			return nil, err
		}
		return r.GetByID(ctx, chat.ID)
	}

	return r.GetByID(ctx, chatID)
}

func (r *GormRepository) AddParticipant(ctx context.Context, chatID uint, userID uint) error {
	p := models.ChatParticipant{
		ChatID: chatID,
		UserID: userID,
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&p).Error
}

func (r *GormRepository) UpdateUnreadCount(ctx context.Context, chatID uint, userID uint, delta int) error {
	if delta > 0 {
		return r.db.WithContext(ctx).Model(&models.ChatParticipant{}).
			Where("chat_id = ? AND user_id = ?", chatID, userID).
			UpdateColumn("unread_count", gorm.Expr("unread_count + ?", delta)).Error
	}
	return nil
}

func (r *GormRepository) ResetUnreadCount(ctx context.Context, chatID uint, userID uint) error {
	return r.db.WithContext(ctx).Model(&models.ChatParticipant{}).
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("unread_count", 0).Error
}

func (r *GormRepository) GetParticipants(ctx context.Context, chatID uint) ([]domain.ChatParticipant, error) {
	var parts []models.ChatParticipant
	if err := r.db.WithContext(ctx).Where("chat_id = ?", chatID).Find(&parts).Error; err != nil {
		return nil, err
	}
	var out []domain.ChatParticipant
	for _, p := range parts {
		out = append(out, domain.ChatParticipant{
			ChatID:            p.ChatID,
			UserID:            p.UserID,
			UnreadCount:       p.UnreadCount,
			LastReadMessageID: p.LastReadMessageID,
			JoinedAt:          p.JoinedAt,
		})
	}
	return out, nil
}
