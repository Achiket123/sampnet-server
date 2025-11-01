package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name           string       `gorm:"type:varchar(255);not null" json:"name"`
	Description    string       `gorm:"type:text" json:"description"`
	OrganisationID uint         `gorm:"index" json:"organisation_id"`
	Organisation   Organisation `gorm:"foreignKey:OrganisationID" json:"organisation"`
	CreatedBy      uint         `gorm:"index" json:"created_by"`
	TeamLead       uint         `gorm:"index" json:"team_lead"`
	TeamLeadUser   UserModel    `gorm:"foreignKey:TeamLead" json:"team_lead_user"`
	CreatedByUser  UserModel    `gorm:"foreignKey:CreatedBy" json:"created_by_user"`
	IsActive       bool         `gorm:"default:true" json:"is_active"`
}

type TeamMember struct {
	gorm.Model
	UserID     uint      `gorm:"index;not null;user_id" json:"user_id"`
	TeamID     uint      `gorm:"index" json:"team_id"`
	User       UserModel `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Team       Team      `gorm:"foreignKey:TeamID;references:ID" json:"team"`
	Role       string    `gorm:"type:varchar(100)" json:"role"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	IsLeader   bool      `gorm:"default:false" json:"is_leader"`
	IsAdmin    bool      `gorm:"default:false" json:"is_admin"`
	IsManager  bool      `gorm:"default:false" json:"is_manager"`
	IsEmployee bool      `gorm:"default:false" json:"is_employee"`
	IsBoss     bool      `gorm:"default:false" json:"is_boss"`
}
