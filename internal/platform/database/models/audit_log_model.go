package models

import (
	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	OrganisationID uint   `gorm:"index;not null"`
	EmployeeID     uint   `gorm:"index;not null"`
	Action         string `gorm:"size:255;not null"` // e.g. "delete_project", "create_employee"
	EntityType     string `gorm:"size:100;not null"` // e.g. "Project", "Employee"
	EntityID       uint   `gorm:"not null"`
	Details        string `gorm:"type:text"`
}