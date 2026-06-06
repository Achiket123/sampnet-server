package resources

import "context"

type RecordFilters struct {
	Search    string                 `json:"search"`
	Filters   map[string]interface{} `json:"filters"`
	SortBy    string                 `json:"sort_by"`
	SortOrder string                 `json:"sort_order"`
	Offset    int                    `json:"offset"`
	Limit     int                    `json:"limit"`
}

type Repository interface {
	CreateCollection(ctx context.Context, collection *ResourceCollection) (*ResourceCollection, error)
	GetCollectionByID(ctx context.Context, id uint, orgID uint) (*ResourceCollection, error)
	GetCollectionsByOrg(ctx context.Context, orgID uint, limit int, offset int) ([]ResourceCollection, int, error)
	UpdateCollection(ctx context.Context, collection *ResourceCollection) (*ResourceCollection, error)
	DeleteCollection(ctx context.Context, id uint, orgID uint) error
	AddField(ctx context.Context, collectionID uint, orgID uint, field FieldDefinition) (*ResourceCollection, error)
	UpdateField(ctx context.Context, collectionID uint, orgID uint, fieldKey string, field FieldDefinition) (*ResourceCollection, error)
	RemoveField(ctx context.Context, collectionID uint, orgID uint, fieldKey string) (*ResourceCollection, error)

	CreateRecord(ctx context.Context, record *ResourceRecord) (*ResourceRecord, error)
	GetRecordByID(ctx context.Context, id uint, collectionID uint, orgID uint) (*ResourceRecord, error)
	GetRecordsByCollection(ctx context.Context, collectionID uint, orgID uint, filters RecordFilters) ([]ResourceRecord, int, error)
	UpdateRecord(ctx context.Context, record *ResourceRecord) (*ResourceRecord, error)
	DeleteRecord(ctx context.Context, id uint, collectionID uint, orgID uint) error
	SearchRecords(ctx context.Context, orgID uint, query string, limit int, offset int) ([]ResourceRecord, error)
}
