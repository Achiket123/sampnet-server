package models

import (
	"time"
)

type TaskComment struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"index;not null"`
	UserID    uint      `gorm:"index;not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	User      UserModel `gorm:"foreignKey:UserID"`
}
