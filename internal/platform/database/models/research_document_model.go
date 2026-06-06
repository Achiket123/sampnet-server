package models

import "gorm.io/gorm"

type ResearchDocument struct {
	gorm.Model
	ResearchID     uint   `gorm:"index;not null" json:"research_id"`
	FolderID       *uint  `gorm:"index" json:"folder_id"`
	OrganisationID uint   `gorm:"index;not null" json:"organisation_id"`
	Title          string `gorm:"size:255;not null" json:"title"`
	Content        string `gorm:"type:text" json:"content"`
	IsPinned       bool   `gorm:"default:false" json:"is_pinned"`
	Status         string `gorm:"size:50;default:'active'" json:"status"`
	CreatedBy      uint   `gorm:"not null" json:"created_by"`
	UpdatedBy      uint   `json:"updated_by"`

	// Relationships
	Versions []ResearchDocumentVersion `gorm:"foreignKey:DocumentID" json:"versions,omitempty"`
}

type ResearchDocumentVersion struct {
	gorm.Model
	DocumentID     uint   `gorm:"index;not null" json:"document_id"`
	OrganisationID uint   `gorm:"index;not null" json:"organisation_id"`
	VersionNumber  int    `gorm:"not null" json:"version_number"`
	Content        string `gorm:"type:text" json:"content"`
	Summary        string `gorm:"size:255" json:"summary"`
	CreatedBy      uint   `gorm:"not null" json:"created_by"`
}
