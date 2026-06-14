package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName string `gorm:"not null" json:"file_name"`
	URL      string `gorm:"type:text;not null" json:"url"`
	FileType string `gorm:"not null" json:"file_type"`
	FileSize int64  `gorm:"not null" json:"file_size"`
}
