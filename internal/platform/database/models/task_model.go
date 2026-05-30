package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title          string        `gorm:"size:255;not null" json:"title"`
	Description    string        `gorm:"type:text" json:"description"`
	DueDate        time.Time     `json:"due_date"`
	AssignedTo     uint          `json:"assigned_to"`
	AssignedUser   *UserModel    `gorm:"foreignKey:AssignedTo;references:ID" json:"assigned_user"`
	AssignedBy     uint          `json:"assigned_by"`
	AssignedByUser *UserModel    `gorm:"foreignKey:AssignedBy;references:ID" json:"assigned_by_user"`
	Type           string        `gorm:"size:100" json:"type"`
	Priority       string        `gorm:"size:100" json:"priority"`
	Status         string        `gorm:"size:50" json:"status"`
	OrganisationID uint          `gorm:"index" json:"organisation_id"`
	Organisation   *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation"`
	IsPersonal     bool          `gorm:"default:false" json:"is_personal"`
	TeamID         *uint         `gorm:"index" json:"team_id"` // Pointer for nullable
	Team           *Team         `gorm:"foreignKey:TeamID" json:"team"`
	ProjectID      *uint         `gorm:"index" json:"project_id"` // Pointer for nullable
	Project        *Project      `gorm:"foreignKey:ProjectID" json:"project"`
}
