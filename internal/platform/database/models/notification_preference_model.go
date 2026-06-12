package models

import (
	"gorm.io/gorm"
)

type NotificationPreference struct {
	gorm.Model
	UserID   uint   `gorm:"index;not null" json:"user_id"`
	Category string `gorm:"size:50;not null" json:"category"` // e.g. "announcements", "tasks", "leaves", "chats"
	Email    bool   `gorm:"default:true" json:"email"`
	Push     bool   `gorm:"default:true" json:"push"`
	InApp    bool   `gorm:"default:true" json:"in_app"`
}
