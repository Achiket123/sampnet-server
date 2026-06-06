package models

import "gorm.io/gorm"

type ResearchFile struct {
	gorm.Model
	ResearchID     uint   `gorm:"index;not null" json:"research_id"`
	DocumentID     *uint  `gorm:"index" json:"document_id"`
	FolderID       *uint  `gorm:"index" json:"folder_id"`
	OrganisationID uint   `gorm:"index;not null" json:"organisation_id"`
	FileName       string `gorm:"size:255;not null" json:"file_name"`
	OriginalName   string `gorm:"size:255;not null" json:"original_name"`
	MimeType       string `gorm:"size:100;not null" json:"mime_type"`
	Extension      string `gorm:"size:20" json:"extension"`
	Size           int64  `gorm:"not null" json:"size"`
	Checksum       string `gorm:"size:255" json:"checksum"`
	StoragePath    string `gorm:"type:text;not null" json:"storage_path"`
	PreviewPath    string `gorm:"type:text" json:"preview_path"`
	ThumbnailPath  string `gorm:"type:text" json:"thumbnail_path"`
	IsPinned       bool   `gorm:"default:false" json:"is_pinned"`
	CreatedBy      uint   `gorm:"not null" json:"created_by"`
	UpdatedBy      uint   `json:"updated_by"`
}
