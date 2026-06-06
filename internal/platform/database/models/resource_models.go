package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CollectionFieldValidation struct {
	MinLength    *int     `json:"min_length"`
	MaxLength    *int     `json:"max_length"`
	MinValue     *float64 `json:"min_value"`
	MaxValue     *float64 `json:"max_value"`
	Pattern      *string  `json:"pattern"`
	ErrorMessage *string  `json:"error_message"`
}

type CollectionRelationConfig struct {
	TargetCollectionID string `json:"target_collection_id"`
	DisplayField       string `json:"display_field"`
	AllowMultiple      bool   `json:"allow_multiple"`
}

type FieldDefinition struct {
	Key               string                     `json:"key"`
	Label             string                     `json:"label"`
	Type              string                     `json:"type"`
	Required          bool                       `json:"required"`
	Unique            bool                       `json:"unique"`
	Options           []string                   `json:"options"`
	DefaultValue      interface{}                `json:"default_value"`
	ValidationRules   *CollectionFieldValidation `json:"validation_rules"`
	FormulaExpression *string                    `json:"formula_expression"`
	RelationConfig    *CollectionRelationConfig  `json:"relation_config"`
	Width             int                        `json:"width"`
	Hidden            bool                       `json:"hidden"`
	Order             int                        `json:"order"`
}

type FieldDefinitions []FieldDefinition

func (f FieldDefinitions) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
	}
	return json.Marshal(f)
}

func (f *FieldDefinitions) Scan(value interface{}) error {
	if value == nil {
		*f = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, f)
}

type RecordData map[string]interface{}

func (r RecordData) Value() (driver.Value, error) {
	if r == nil {
		return nil, nil
	}
	return json.Marshal(r)
}

func (r *RecordData) Scan(value interface{}) error {
	if value == nil {
		*r = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, r)
}

type JSONMapResource map[string]interface{}

func (m JSONMapResource) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(m)
}

func (m *JSONMapResource) Scan(value interface{}) error {
	if value == nil {
		*m = map[string]interface{}{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}

type JSONListResource []map[string]interface{}

func (m JSONListResource) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal([]map[string]interface{}{})
	}
	return json.Marshal(m)
}

func (m *JSONListResource) Scan(value interface{}) error {
	if value == nil {
		*m = []map[string]interface{}{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}


type ResourceCollection struct {
	gorm.Model
	OrganisationID uint             `gorm:"index;not null" json:"organisation_id"`
	Name           string           `gorm:"size:255;not null" json:"name"`
	Description    string           `gorm:"type:text" json:"description"`
	Icon           *string          `gorm:"size:100" json:"icon"`
	Colour         *string          `gorm:"size:100;default:'#6B7280'" json:"colour"`
	DefaultView    string           `gorm:"size:20;default:'table'" json:"default_view"`
	KanbanField    *string          `gorm:"size:100" json:"kanban_field"`
	CreatedBy      uint             `json:"created_by"`
	Fields         FieldDefinitions `gorm:"type:jsonb" json:"fields"`
	SortConfig     JSONListResource `gorm:"type:jsonb;default:'[]'" json:"sort_config"`
	FilterConfig   JSONListResource `gorm:"type:jsonb;default:'[]'" json:"filter_config"`
}

type ResourceRecord struct {
	gorm.Model
	CollectionID   uint         `gorm:"index;not null" json:"collection_id"`
	OrganisationID uint         `gorm:"index;not null" json:"organisation_id"`
	Version        int          `gorm:"default:1" json:"version"`
	Data           RecordData   `gorm:"type:jsonb" json:"data"`
	CreatedBy      uint         `json:"created_by"`
	UpdatedBy      *uint        `json:"updated_by"`
	ProjectID      *uint        `gorm:"index" json:"project_id"`
	TeamID         *uint        `gorm:"index" json:"team_id"`
	TaskID         *uint        `gorm:"index" json:"task_id"`
}

type ResourceRecordHistory struct {
	gorm.Model
	RecordID     uint            `gorm:"index;not null"`
	CollectionID uint            `gorm:"index;not null"`
	ChangedByID  uint            `gorm:"index;not null"`
	Changes      JSONMapResource `gorm:"type:jsonb;not null"`
	Version      int             `gorm:"not null"`
	ChangedAt    time.Time       `gorm:"default:current_timestamp"`
	ChangedBy    Employee        `gorm:"foreignKey:ChangedByID;references:UserID"`
}

type ResourceRecordAttachment struct {
	gorm.Model
	RecordID     uint      `gorm:"index;not null"`
	CollectionID uint      `gorm:"index;not null"`
	FieldKey     string    `gorm:"size:100;not null"`
	FileName     string    `gorm:"size:500;not null"`
	FileSize     int64     `gorm:"not null"`
	MimeType     string    `gorm:"size:200"`
	FileUrl      string    `gorm:"type:text;not null"`
	UploadedByID uint      `gorm:"index;not null"`
	UploadedAt   time.Time `gorm:"default:current_timestamp"`
}
