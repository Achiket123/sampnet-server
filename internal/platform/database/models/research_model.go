package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type TagsList []string

func (t TagsList) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}

func (t *TagsList) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, t)
}

type ResearchEntry struct {
	gorm.Model
	Title          string       `gorm:"size:255;not null" json:"title"`
	Description    string       `gorm:"type:text" json:"description"`
	Thumbnail      string       `gorm:"type:varchar(512)" json:"thumbnail"`
	Status         string       `gorm:"size:50;not null;default:'active'" json:"status"` // draft, active, archived, completed
	Visibility     string       `gorm:"size:50;not null;default:'private'" json:"visibility"` // private, team, organisation
	CreatedBy      uint         `gorm:"not null" json:"created_by"`
	UpdatedBy      uint         `json:"updated_by"`
	Author         UserModel    `gorm:"foreignKey:CreatedBy" json:"author"`
	ProjectID      *uint        `gorm:"index" json:"project_id"`
	Project        Project      `gorm:"foreignKey:ProjectID" json:"project"`
	TeamID         *uint        `gorm:"index" json:"team_id"`
	Team           Team         `gorm:"foreignKey:TeamID" json:"team"`
	OrganisationID uint         `gorm:"index;not null" json:"organisation_id"`
	Tags           TagsList     `gorm:"type:jsonb" json:"tags"`
	
	// Relationships
	Folders       []ResearchFolder       `gorm:"foreignKey:ResearchID" json:"folders"`
	Documents     []ResearchDocument     `gorm:"foreignKey:ResearchID" json:"documents"`
	Files         []ResearchFile         `gorm:"foreignKey:ResearchID" json:"files"`
	Collaborators []ResearchCollaborator `gorm:"foreignKey:ResearchID" json:"collaborators"`
}
