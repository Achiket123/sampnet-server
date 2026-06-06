package models

import (
	"time"

	"gorm.io/gorm"
)

type Leave struct {
	gorm.Model
	EmployeeID     uint      `gorm:"index;not null"`
	OrganisationID uint      `gorm:"index;not null"`
	ManagerID      *uint     `gorm:"index"`
	LeaveType      string    `gorm:"size:50;not null"`
	StartDate      time.Time `gorm:"not null"`
	EndDate        time.Time `gorm:"not null"`
	TotalDays      int       `gorm:"not null"`
	Reason         string    `gorm:"type:text"`
	Status         string    `gorm:"size:50;not null;default:'pending'"`
	ManagerNote    string    `gorm:"type:text"`

	Employee UserModel `gorm:"foreignKey:EmployeeID"`
	Manager  *UserModel `gorm:"foreignKey:ManagerID"`
}
