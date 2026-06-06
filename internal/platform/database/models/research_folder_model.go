package models

import "gorm.io/gorm"

type ResearchFolder struct {
	gorm.Model
	ResearchID     uint            `gorm:"index;not null" json:"research_id"`
	OrganisationID uint            `gorm:"index;not null" json:"organisation_id"`
	ParentID       *uint           `gorm:"index" json:"parent_id"`
	Parent         *ResearchFolder `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Name           string          `gorm:"size:255;not null" json:"name"`
	CreatedBy      uint            `gorm:"not null" json:"created_by"`
	UpdatedBy      uint            `json:"updated_by"`

	// Relationships
	SubFolders []ResearchFolder   `gorm:"foreignKey:ParentID" json:"sub_folders,omitempty"`
	Documents  []ResearchDocument `gorm:"foreignKey:FolderID" json:"documents,omitempty"`
	Files      []ResearchFile     `gorm:"foreignKey:FolderID" json:"files,omitempty"`
}
