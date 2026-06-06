package models

import (
	"gorm.io/gorm"
)

type TaskType struct {
	gorm.Model
	OrganisationID uint   `gorm:"index;not null"`
	Name           string `gorm:"size:100;not null"`
	Description    string `gorm:"type:text"`
}