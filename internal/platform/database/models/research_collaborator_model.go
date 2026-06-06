package models

import "gorm.io/gorm"

type ResearchCollaborator struct {
	gorm.Model
	ResearchID uint      `gorm:"index;not null" json:"research_id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	User       UserModel `gorm:"foreignKey:UserID" json:"user"`
	Role       string    `gorm:"size:50;not null;default:'viewer'" json:"role"` // owner, editor, commenter, viewer
	AddedBy    uint      `json:"added_by"`
}
