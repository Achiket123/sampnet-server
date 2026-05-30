package models

import (
	"time"
)

type TaskAttachment struct {
	ID         uint      `gorm:"primaryKey"`
	TaskID     uint      `gorm:"index;not null"`
	FileID     uint      `gorm:"index;not null"`
	UploadedBy uint      `gorm:"index;not null"`
	FileName   string    `gorm:"size:255"`
	CreatedAt  time.Time
	File       File      `gorm:"foreignKey:FileID"`
	UploadedByUser UserModel `gorm:"foreignKey:UploadedBy"`
}
