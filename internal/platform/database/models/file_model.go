package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName string `gorm:"not null" json:"file_name"`
	Data     []byte `gorm:"type:bytea;not null" json:"data"`
	FileType string `gorm:"not null" json:"file_type"`
	FileSize int64  `gorm:"not null" json:"file_size"`
}
