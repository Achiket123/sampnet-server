package models

import "gorm.io/gorm"

type ResearchActivity struct {
	gorm.Model
	ResearchID uint   `gorm:"index;not null" json:"research_id"`
	ActorID    uint   `gorm:"index;not null" json:"actor_id"`
	Action     string `gorm:"size:100;not null" json:"action"` // e.g., document_created, folder_deleted
	EntityType string `gorm:"size:100;not null" json:"entity_type"` // e.g., document, file, folder
	EntityID   uint   `json:"entity_id"`
	Metadata   string `gorm:"type:jsonb" json:"metadata"`
}
