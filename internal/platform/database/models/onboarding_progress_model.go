package models

import "time"

type OnboardingProgress struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"index;not null" json:"userId"`
	OrganisationID    uint      `gorm:"index" json:"organisationId"`
	ProfileCompleted  bool      `gorm:"default:false" json:"profileCompleted"`
	TeamJoined         bool      `gorm:"default:false" json:"teamJoined"`
	TaskCreated       bool      `gorm:"default:false" json:"taskCreated"`
	InviteSent        bool      `gorm:"default:false" json:"inviteSent"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
