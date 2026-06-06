package models

import (
	"time"

	"gorm.io/gorm"
)

type Milestone struct {
	gorm.Model
	ProjectID      uint      `gorm:"index" json:"project_id"`
	Project        Project   `gorm:"foreignKey:ProjectID" json:"project"`
	Title          string    `gorm:"size:255" json:"title"`
	Description    string    `gorm:"type:text" json:"description"`
	DueDate        time.Time `json:"due_date"`
	Status         string    `gorm:"size:50" json:"status"` // e.g. "Pending", "Completed"
	OrganisationID uint      `gorm:"index" json:"organisation_id"`
}