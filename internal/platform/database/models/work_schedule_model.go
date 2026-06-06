package models

import (
	"gorm.io/gorm"
)

type WorkSchedule struct {
	gorm.Model
	OrganisationID uint   `gorm:"index;not null"`
	Name           string `gorm:"size:100;not null"` // e.g. "Standard Week", "Flexible Shift"
	ScheduleData   string `gorm:"type:text"`         // JSON or text describing work hours per day
}