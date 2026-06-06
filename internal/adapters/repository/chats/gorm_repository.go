package chats

import (
	"context"
	domain "server/internal/domain/chats"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository { return &GormRepository{db: db} }

func (r *GormRepository) Create(ctx context.Context, chat *domain.Chat) error {
	m := models.Chat{ID: chat.ID, FirstName: chat.FirstName, LastName: chat.LastName, Email: chat.Email, OrganisationID: chat.OrganisationID, LastMessage: chat.LastMessage, LastMessageTimestamp: chat.LastMessageTimestamp, NumberOfMessage: chat.NumberOfMessage}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&m).Error
}

func (r *GormRepository) ListByOrganisation(ctx context.Context, organisationID uint) ([]domain.Chat, error) {
	// Automatically synchronize/upsert missing employees/managers/boss of that organisation into the chats table.
	var employees []models.Employee
	_ = r.db.WithContext(ctx).Preload("User").Where("organisation_id = ?", organisationID).Find(&employees)

	var managers []models.Manager
	_ = r.db.WithContext(ctx).Preload("User").Where("organisation_id = ?", organisationID).Find(&managers)

	var boss models.Boss
	bossFound := false
	if err := r.db.WithContext(ctx).Preload("User").Where("organisation_id = ?", organisationID).First(&boss).Error; err == nil {
		bossFound = true
	}

	now := time.Now()

	for _, emp := range employees {
		if emp.User.ID > 0 {
			var count int64
			r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", emp.User.ID).Count(&count)
			if count == 0 {
				chat := models.Chat{
					ID:                   emp.User.ID,
					FirstName:            emp.User.FirstName,
					LastName:             emp.User.LastName,
					Email:                emp.User.Email,
					OrganisationID:       organisationID,
					LastMessage:          "",
					LastMessageTimestamp: &now,
					NumberOfMessage:      0,
				}
				r.db.WithContext(ctx).Create(&chat)
			}
		}
	}

	for _, mgr := range managers {
		if mgr.User.ID > 0 {
			var count int64
			r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", mgr.User.ID).Count(&count)
			if count == 0 {
				chat := models.Chat{
					ID:                   mgr.User.ID,
					FirstName:            mgr.User.FirstName,
					LastName:             mgr.User.LastName,
					Email:                mgr.User.Email,
					OrganisationID:       organisationID,
					LastMessage:          "",
					LastMessageTimestamp: &now,
					NumberOfMessage:      0,
				}
				r.db.WithContext(ctx).Create(&chat)
			}
		}
	}

	if bossFound && boss.User.ID > 0 {
		var count int64
		r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", boss.User.ID).Count(&count)
		if count == 0 {
			chat := models.Chat{
				ID:                   boss.User.ID,
				FirstName:            boss.User.FirstName,
				LastName:             boss.User.LastName,
				Email:                boss.User.Email,
				OrganisationID:       organisationID,
				LastMessage:          "",
				LastMessageTimestamp: &now,
				NumberOfMessage:      0,
			}
			r.db.WithContext(ctx).Create(&chat)
		}
	}

	var rows []models.Chat
	if err := r.db.WithContext(ctx).Where("organisation_id = ?", organisationID).Order("last_message_timestamp desc nulls last").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Chat, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Chat{ID: row.ID, FirstName: row.FirstName, LastName: row.LastName, Email: row.Email, OrganisationID: row.OrganisationID, LastMessage: row.LastMessage, LastMessageTimestamp: row.LastMessageTimestamp, NumberOfMessage: row.NumberOfMessage})
	}
	return out, nil
}

func (r *GormRepository) UpdateLastMessage(ctx context.Context, id uint, message string, at *string) error {
	var ts *time.Time
	if at != nil && *at != "" {
		parsed, err := time.Parse(time.RFC3339, *at)
		if err == nil {
			ts = &parsed
		}
	}
	updates := map[string]any{"last_message": message, "last_message_timestamp": ts}
	return r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", id).Updates(updates).Error
}

func (r *GormRepository) IncrementMessageCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Chat{}).Where("id = ?", id).UpdateColumn("number_of_message", gorm.Expr("number_of_message + 1")).Error
}
