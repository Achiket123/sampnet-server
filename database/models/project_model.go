package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name             string
	Description      string
	StartDate        time.Time
	EndDate          time.Time
	OrganisationID   uint         `gorm:"index" json:"organisation_id"`
	TeamID           uint         `gorm:"index" json:"team_id"`
	Team             Team         `gorm:"foreignKey:TeamID" json:"team"`
	CreatedBy        uint         `gorm:"index" json:"created_by"`
	CreatedByUser    UserModel    `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	Status           string       `gorm:"size:50" json:"status"`
	Priority         string       `gorm:"size:50" json:"priority"`
	CompletionStatus string       `gorm:"size:50" json:"completion_status"`
	Organisation     Organisation `gorm:"foreignKey:OrganisationID" json:"organisation"`
}
