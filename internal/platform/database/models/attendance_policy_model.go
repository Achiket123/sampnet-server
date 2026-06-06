package models

import (
	"gorm.io/gorm"
)

type AttendancePolicy struct {
	gorm.Model
	OrganisationID uint   `gorm:"index;not null"`
	CheckInTime    string `gorm:"size:10;not null"` // e.g., "09:00"
	CheckOutTime   string `gorm:"size:10;not null"` // e.g., "18:00"
	GracePeriodMins int   `gorm:"default:15"`
}