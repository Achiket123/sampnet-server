package models

import "gorm.io/gorm"

type Notification struct {

	gorm.Model
	UserID uint `gorm:"index" json:"user_id"`
	OrganisationID uint `gorm:"index" json:"organisation_id"`
	User           UserModel      `gorm:"foreignKey:UserID" json:"user"`
	Organisation   Organisation  `gorm:"foreignKey:OrganisationID" json:"organisation"`
	Title          string        `json:"title"`
	Message        string        `json:"message"`
	IsRead         bool          `json:"is_read"`
	NotificationType string        `json:"notification_type"`
	Link             string        `json:"link"`
}
