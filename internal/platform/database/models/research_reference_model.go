package models

import "gorm.io/gorm"

type ResearchReference struct {
	gorm.Model
	ResearchID     uint   `gorm:"index;not null" json:"research_id"`
	DocumentID     *uint  `gorm:"index" json:"document_id"`
	OrganisationID uint   `gorm:"index;not null" json:"organisation_id"`
	Title          string `gorm:"size:255;not null" json:"title"`
	Authors        string `gorm:"size:512" json:"authors"`
	Publisher      string `gorm:"size:255" json:"publisher"`
	PublicationDate string `gorm:"size:100" json:"publication_date"`
	DOI            string `gorm:"size:100" json:"doi"`
	URL            string `gorm:"type:text" json:"url"`
	Journal        string `gorm:"size:255" json:"journal"`
	Conference     string `gorm:"size:255" json:"conference"`
	Citation       string `gorm:"type:text" json:"citation"`
	Notes          string `gorm:"type:text" json:"notes"`
	CreatedBy      uint   `gorm:"not null" json:"created_by"`
}
