package models

import (
	"gorm.io/gorm"
)

type RolePermissions struct {
	gorm.Model
	OrganisationID uint   `gorm:"index;not null"`
	Role           string `gorm:"size:50;not null"` // e.g. "employee", "manager", "admin"
	Permission     string `gorm:"size:100;not null"` // e.g. "read:employees", "write:tasks"
}