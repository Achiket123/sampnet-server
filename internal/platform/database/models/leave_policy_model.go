package models

import (
	"gorm.io/gorm"
)

type LeavePolicy struct {
	gorm.Model
	OrganisationID uint   `gorm:"index;not null"`
	LeaveType      string `gorm:"size:50;not null"`
	MaxDays        int    `gorm:"not null"`
	Description    string `gorm:"type:text"`
}