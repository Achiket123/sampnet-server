package models

import (
	"gorm.io/gorm"
)

type TaskActivity struct {
	gorm.Model
	TaskID       uint   `gorm:"index;not null"`
	EmployeeID   uint   `gorm:"index;not null"`
	ActivityType string `gorm:"size:50;not null"` // e.g. "created", "status_changed", "assigned"
	Description  string `gorm:"type:text"`
}