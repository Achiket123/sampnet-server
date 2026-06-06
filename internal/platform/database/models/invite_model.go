package models

import (
	"time"

	"gorm.io/gorm"
)

type Invite struct {
	gorm.Model
	OrganisationID uint      `gorm:"index;not null"`
	Email          string    `gorm:"size:255;not null"`
	Role           string    `gorm:"size:50;not null"` // e.g. "employee", "manager"
	Token          string    `gorm:"size:255;not null;unique"`
	ExpiresAt      time.Time `gorm:"not null"`
	Status         string    `gorm:"size:50;not null;default:'pending'"` // e.g. pending, accepted, expired
}