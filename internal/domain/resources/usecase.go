package resources

import "context"

type BulkCreateResult struct {
	Created int                      `json:"created"`
	Failed  int                      `json:"failed"`
	Errors  []BulkCreateRecordError  `json:"errors"`
}

type BulkCreateRecordError struct {
	Row     int               `json:"row"`
	Field   string            `json:"field"`
	Message string            `json:"message"`
}

type UseCase interface {
	CreateCollection(ctx context.Context, orgID uint, createdBy uint, name string, description string, icon *string, colour *string, fields []FieldDefinition) (*ResourceCollection, error)
	GetCollection(ctx context.Context, id uint, orgID uint) (*ResourceCollection, error)
	ListCollections(ctx context.Context, orgID uint, limit int, offset int) ([]ResourceCollection, int, error)
	UpdateCollection(ctx context.Context, id uint, orgID uint, name string, description string, icon *string, colour *string) (*ResourceCollection, error)
	DeleteCollection(ctx context.Context, id uint, orgID uint, force bool) error
	AddFieldToCollection(ctx context.Context, collectionID uint, orgID uint, field FieldDefinition) (*ResourceCollection, error)
	UpdateFieldInCollection(ctx context.Context, collectionID uint, orgID uint, fieldKey string, field FieldDefinition) (*ResourceCollection, error)
	RemoveFieldFromCollection(ctx context.Context, collectionID uint, orgID uint, fieldKey string) (*ResourceCollection, string, error)

	CreateRecord(ctx context.Context, orgID uint, collectionID uint, createdBy uint, data map[string]interface{}, projectID *uint, teamID *uint, taskID *uint) (*ResourceRecord, error)
	GetRecord(ctx context.Context, id uint, collectionID uint, orgID uint) (*ResourceRecord, error)
	ListRecords(ctx context.Context, collectionID uint, orgID uint, filters RecordFilters) ([]ResourceRecord, int, error)
	UpdateRecord(ctx context.Context, id uint, collectionID uint, orgID uint, updatedBy uint, data map[string]interface{}) (*ResourceRecord, error)
	DeleteRecord(ctx context.Context, id uint, collectionID uint, orgID uint) error
	BulkCreateRecords(ctx context.Context, orgID uint, collectionID uint, createdBy uint, records []map[string]interface{}) (*BulkCreateResult, error)
	ExportRecords(ctx context.Context, collectionID uint, orgID uint) ([]map[string]interface{}, []string, error)
}
