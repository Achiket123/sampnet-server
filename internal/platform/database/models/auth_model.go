package models

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt int64     `gorm:"not null"`
	Revoked   bool      `gorm:"not null"`
	User      UserModel `gorm:"foreignKey:UserID"`
}
