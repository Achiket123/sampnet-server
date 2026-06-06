package resources

import (
	"fmt"
	"time"
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

type ResourceCollection struct {
	ID             uint                   `json:"id"`
	OrganisationID uint                   `json:"organisation_id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Icon           *string                `json:"icon"`
	Colour         *string                `json:"colour"`
	DefaultView    string                 `json:"default_view"`
	KanbanField    *string                `json:"kanban_field"`
	CreatedBy      uint                   `json:"created_by"`
	Fields         []FieldDefinition      `json:"fields"`
	SortConfig     []map[string]interface{} `json:"sort_config"`
	FilterConfig   []map[string]interface{} `json:"filter_config"`
	RecordCount    int                    `json:"record_count"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type ResourceRecord struct {
	ID             uint                   `json:"id"`
	CollectionID   uint                   `json:"collection_id"`
	OrganisationID uint                   `json:"organisation_id"`
	Version        int                    `json:"version"`
	Data           map[string]interface{} `json:"data"`
	CreatedBy      uint                   `json:"created_by"`
	CreatedByName  *string                `json:"created_by_name"`
	UpdatedBy      *uint                  `json:"updated_by"`
	UpdatedByName  *string                `json:"updated_by_name"`
	ProjectID      *uint                  `json:"project_id"`
	TeamID         *uint                  `json:"team_id"`
	TaskID         *uint                  `json:"task_id"`
	Attachments    []RecordAttachment     `json:"attachments"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type RecordAttachment struct {
	ID           uint      `json:"id"`
	RecordID     uint      `json:"record_id"`
	CollectionID uint      `json:"collection_id"`
	FieldKey     string    `json:"field_key"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	FileUrl      string    `json:"file_url"`
	UploadedByID uint      `json:"uploaded_by_id"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

type RecordHistory struct {
	ID              uint                   `json:"id"`
	RecordID        uint                   `json:"record_id"`
	CollectionID    uint                   `json:"collection_id"`
	ChangedByID     uint                   `json:"changed_by_id"`
	ChangedByName   string                 `json:"changed_by_name"`
	ChangedByAvatar *string                `json:"changed_by_avatar"`
	Changes         map[string]interface{} `json:"changes"`
	Version         int                    `json:"version"`
	ChangedAt       time.Time              `json:"changed_at"`
}

type RecordsPage struct {
	Records []ResourceRecord `json:"records"`
	Total   int              `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
	HasMore bool             `json:"has_more"`
}

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", v.Errors)
}
