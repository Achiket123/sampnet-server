package resources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"server/internal/domain/resources"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) resources.Repository {
	return &gormRepository{db: db}
}

func toDomainCollection(m *models.ResourceCollection) *resources.ResourceCollection {
	fields := make([]resources.FieldDefinition, len(m.Fields))
	for i, f := range m.Fields {
		fields[i] = resources.FieldDefinition{
			Key:          f.Key,
			Label:        f.Label,
			Type:         f.Type,
			Required:          f.Required,
			Unique:            f.Unique,
			Options:           f.Options,
			DefaultValue:      f.DefaultValue,
			FormulaExpression: f.FormulaExpression,
			Width:             f.Width,
			Hidden:            f.Hidden,
			Order:             f.Order,
		}
		if f.ValidationRules != nil {
			fields[i].ValidationRules = &resources.CollectionFieldValidation{
				MinLength:    f.ValidationRules.MinLength,
				MaxLength:    f.ValidationRules.MaxLength,
				MinValue:     f.ValidationRules.MinValue,
				MaxValue:     f.ValidationRules.MaxValue,
				Pattern:      f.ValidationRules.Pattern,
				ErrorMessage: f.ValidationRules.ErrorMessage,
			}
		}
		if f.RelationConfig != nil {
			fields[i].RelationConfig = &resources.CollectionRelationConfig{
				TargetCollectionID: f.RelationConfig.TargetCollectionID,
				DisplayField:       f.RelationConfig.DisplayField,
				AllowMultiple:      f.RelationConfig.AllowMultiple,
			}
		}
	}
	return &resources.ResourceCollection{
		ID:             m.ID,
		OrganisationID: m.OrganisationID,
		Name:           m.Name,
		Description:    m.Description,
		Icon:           m.Icon,
		Colour:         m.Colour,
		DefaultView:    m.DefaultView,
		KanbanField:    m.KanbanField,
		CreatedBy:      m.CreatedBy,
		Fields:         fields,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func toDomainRecord(m *models.ResourceRecord) *resources.ResourceRecord {
	return &resources.ResourceRecord{
		ID:             m.ID,
		CollectionID:   m.CollectionID,
		OrganisationID: m.OrganisationID,
		Version:        m.Version,
		Data:           map[string]interface{}(m.Data),
		CreatedBy:      m.CreatedBy,
		UpdatedBy:      m.UpdatedBy,
		ProjectID:      m.ProjectID,
		TeamID:         m.TeamID,
		TaskID:         m.TaskID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (r *gormRepository) CreateCollection(ctx context.Context, collection *resources.ResourceCollection) (*resources.ResourceCollection, error) {
	fields := make(models.FieldDefinitions, len(collection.Fields))
	for i, f := range collection.Fields {
		fields[i] = models.FieldDefinition{
			Key:          f.Key,
			Label:        f.Label,
			Type:         f.Type,
			Required:          f.Required,
			Unique:            f.Unique,
			Options:           f.Options,
			DefaultValue:      f.DefaultValue,
			FormulaExpression: f.FormulaExpression,
			Width:             f.Width,
			Hidden:            f.Hidden,
			Order:             f.Order,
		}
		if f.ValidationRules != nil {
			fields[i].ValidationRules = &models.CollectionFieldValidation{
				MinLength:    f.ValidationRules.MinLength,
				MaxLength:    f.ValidationRules.MaxLength,
				MinValue:     f.ValidationRules.MinValue,
				MaxValue:     f.ValidationRules.MaxValue,
				Pattern:      f.ValidationRules.Pattern,
				ErrorMessage: f.ValidationRules.ErrorMessage,
			}
		}
		if f.RelationConfig != nil {
			fields[i].RelationConfig = &models.CollectionRelationConfig{
				TargetCollectionID: f.RelationConfig.TargetCollectionID,
				DisplayField:       f.RelationConfig.DisplayField,
				AllowMultiple:      f.RelationConfig.AllowMultiple,
			}
		}
	}
	dbCollection := models.ResourceCollection{
		OrganisationID: collection.OrganisationID,
		Name:           collection.Name,
		Description:    collection.Description,
		Icon:           collection.Icon,
		Colour:         collection.Colour,
		DefaultView:    collection.DefaultView,
		KanbanField:    collection.KanbanField,
		CreatedBy:      collection.CreatedBy,
		Fields:         fields,
	}
	if err := r.db.WithContext(ctx).Create(&dbCollection).Error; err != nil {
		return nil, err
	}
	return toDomainCollection(&dbCollection), nil
}

func (r *gormRepository) GetCollectionByID(ctx context.Context, id uint, orgID uint) (*resources.ResourceCollection, error) {
	var collection models.ResourceCollection
	if err := r.db.WithContext(ctx).Where("id = ? AND organisation_id = ?", id, orgID).First(&collection).Error; err != nil {
		return nil, err
	}
	return toDomainCollection(&collection), nil
}

func (r *gormRepository) GetCollectionsByOrg(ctx context.Context, orgID uint, limit int, offset int) ([]resources.ResourceCollection, int, error) {
	var collections []models.ResourceCollection
	var count int64

	tx := r.db.WithContext(ctx).Model(&models.ResourceCollection{}).Where("organisation_id = ?", orgID)
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := tx.Limit(limit).Offset(offset).Find(&collections).Error; err != nil {
		return nil, 0, err
	}

	domainColls := make([]resources.ResourceCollection, len(collections))
	for i, coll := range collections {
		domainColls[i] = *toDomainCollection(&coll)
	}

	return domainColls, int(count), nil
}

func (r *gormRepository) UpdateCollection(ctx context.Context, collection *resources.ResourceCollection) (*resources.ResourceCollection, error) {
	var dbCollection models.ResourceCollection
	if err := r.db.WithContext(ctx).Where("id = ? AND organisation_id = ?", collection.ID, collection.OrganisationID).First(&dbCollection).Error; err != nil {
		return nil, err
	}

	dbCollection.Name = collection.Name
	dbCollection.Description = collection.Description
	dbCollection.Icon = collection.Icon
	dbCollection.Colour = collection.Colour
	dbCollection.DefaultView = collection.DefaultView
	dbCollection.KanbanField = collection.KanbanField

	if err := r.db.WithContext(ctx).Save(&dbCollection).Error; err != nil {
		return nil, err
	}

	return toDomainCollection(&dbCollection), nil
}

func (r *gormRepository) DeleteCollection(ctx context.Context, id uint, orgID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var collection models.ResourceCollection
		if err := tx.Where("id = ? AND organisation_id = ?", id, orgID).First(&collection).Error; err != nil {
			return err
		}
		// Soft delete all records in collection
		if err := tx.Where("collection_id = ? AND organisation_id = ?", id, orgID).Delete(&models.ResourceRecord{}).Error; err != nil {
			return err
		}
		// Soft delete collection
		return tx.Delete(&collection).Error
	})
}

func (r *gormRepository) AddField(ctx context.Context, collectionID uint, orgID uint, field resources.FieldDefinition) (*resources.ResourceCollection, error) {
	var collection models.ResourceCollection
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organisation_id = ?", collectionID, orgID).First(&collection).Error; err != nil {
			return err
		}

		for _, f := range collection.Fields {
			if f.Key == field.Key {
				return errors.New("field key already exists")
			}
		}

		newField := models.FieldDefinition{
			Key:          field.Key,
			Label:        field.Label,
			Type:         field.Type,
			Required:          field.Required,
			Unique:            field.Unique,
			Options:           field.Options,
			DefaultValue:      field.DefaultValue,
			FormulaExpression: field.FormulaExpression,
			Width:             field.Width,
			Hidden:            field.Hidden,
			Order:             field.Order,
		}
		if field.ValidationRules != nil {
			newField.ValidationRules = &models.CollectionFieldValidation{
				MinLength:    field.ValidationRules.MinLength,
				MaxLength:    field.ValidationRules.MaxLength,
				MinValue:     field.ValidationRules.MinValue,
				MaxValue:     field.ValidationRules.MaxValue,
				Pattern:      field.ValidationRules.Pattern,
				ErrorMessage: field.ValidationRules.ErrorMessage,
			}
		}
		if field.RelationConfig != nil {
			newField.RelationConfig = &models.CollectionRelationConfig{
				TargetCollectionID: field.RelationConfig.TargetCollectionID,
				DisplayField:       field.RelationConfig.DisplayField,
				AllowMultiple:      field.RelationConfig.AllowMultiple,
			}
		}
		collection.Fields = append(collection.Fields, newField)

		return tx.Save(&collection).Error
	})
	if err != nil {
		return nil, err
	}
	return toDomainCollection(&collection), nil
}

func (r *gormRepository) UpdateField(ctx context.Context, collectionID uint, orgID uint, fieldKey string, field resources.FieldDefinition) (*resources.ResourceCollection, error) {
	var collection models.ResourceCollection
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organisation_id = ?", collectionID, orgID).First(&collection).Error; err != nil {
			return err
		}

		found := false
		for i, f := range collection.Fields {
			if f.Key == fieldKey {
				// Type changes are validated at UseCase layer.
				collection.Fields[i].Label = field.Label
				collection.Fields[i].Required = field.Required
				collection.Fields[i].Unique = field.Unique
				collection.Fields[i].Options = field.Options
				collection.Fields[i].DefaultValue = field.DefaultValue
				collection.Fields[i].FormulaExpression = field.FormulaExpression
				collection.Fields[i].Width = field.Width
				collection.Fields[i].Hidden = field.Hidden
				collection.Fields[i].Order = field.Order
				
				if field.ValidationRules != nil {
					collection.Fields[i].ValidationRules = &models.CollectionFieldValidation{
						MinLength:    field.ValidationRules.MinLength,
						MaxLength:    field.ValidationRules.MaxLength,
						MinValue:     field.ValidationRules.MinValue,
						MaxValue:     field.ValidationRules.MaxValue,
						Pattern:      field.ValidationRules.Pattern,
						ErrorMessage: field.ValidationRules.ErrorMessage,
					}
				} else {
					collection.Fields[i].ValidationRules = nil
				}
				
				if field.RelationConfig != nil {
					collection.Fields[i].RelationConfig = &models.CollectionRelationConfig{
						TargetCollectionID: field.RelationConfig.TargetCollectionID,
						DisplayField:       field.RelationConfig.DisplayField,
						AllowMultiple:      field.RelationConfig.AllowMultiple,
					}
				} else {
					collection.Fields[i].RelationConfig = nil
				}
				
				found = true
				break
			}
		}

		if !found {
			return errors.New("field key not found")
		}

		return tx.Save(&collection).Error
	})
	if err != nil {
		return nil, err
	}
	return toDomainCollection(&collection), nil
}

func (r *gormRepository) RemoveField(ctx context.Context, collectionID uint, orgID uint, fieldKey string) (*resources.ResourceCollection, error) {
	var collection models.ResourceCollection
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organisation_id = ?", collectionID, orgID).First(&collection).Error; err != nil {
			return err
		}

		found := false
		var newFields models.FieldDefinitions
		for _, f := range collection.Fields {
			if f.Key == fieldKey {
				found = true
			} else {
				newFields = append(newFields, f)
			}
		}

		if !found {
			return errors.New("field key not found")
		}

		collection.Fields = newFields
		return tx.Save(&collection).Error
	})
	if err != nil {
		return nil, err
	}
	return toDomainCollection(&collection), nil
}

func (r *gormRepository) CreateRecord(ctx context.Context, record *resources.ResourceRecord) (*resources.ResourceRecord, error) {
	dbRecord := models.ResourceRecord{
		CollectionID:   record.CollectionID,
		OrganisationID: record.OrganisationID,
		Data:           models.RecordData(record.Data),
		CreatedBy:      record.CreatedBy,
		ProjectID:      record.ProjectID,
		TeamID:         record.TeamID,
		TaskID:         record.TaskID,
	}
	if err := r.db.WithContext(ctx).Create(&dbRecord).Error; err != nil {
		return nil, err
	}
	return toDomainRecord(&dbRecord), nil
}

func (r *gormRepository) GetRecordByID(ctx context.Context, id uint, collectionID uint, orgID uint) (*resources.ResourceRecord, error) {
	var rec models.ResourceRecord
	if err := r.db.WithContext(ctx).Where("id = ? AND collection_id = ? AND organisation_id = ?", id, collectionID, orgID).First(&rec).Error; err != nil {
		return nil, err
	}
	return toDomainRecord(&rec), nil
}

func (r *gormRepository) GetRecordsByCollection(ctx context.Context, collectionID uint, orgID uint, filters resources.RecordFilters) ([]resources.ResourceRecord, int, error) {
	var recs []models.ResourceRecord
	var count int64

	tx := r.db.WithContext(ctx).Model(&models.ResourceRecord{}).Where("collection_id = ? AND organisation_id = ?", collectionID, orgID)

	// JSONB full text cast search
	if filters.Search != "" {
		tx = tx.Where("data::text ILIKE ?", "%"+filters.Search+"%")
	}

	// JSONB containment filter
	if len(filters.Filters) > 0 {
		filterJSON, err := json.Marshal(filters.Filters)
		if err == nil {
			tx = tx.Where("data @> ?::jsonb", filterJSON)
		}
	}

	// Count total
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	if filters.SortBy != "" {
		order := "ASC"
		if strings.ToUpper(filters.SortOrder) == "DESC" {
			order = "DESC"
		}
		if filters.SortBy == "created_at" || filters.SortBy == "updated_at" || filters.SortBy == "id" {
			tx = tx.Order(fmt.Sprintf("%s %s", filters.SortBy, order))
		} else {
			sanitizedField := sanitizeFieldName(filters.SortBy)
			tx = tx.Order(fmt.Sprintf("data->>'%s' %s", sanitizedField, order))
		}
	} else {
		tx = tx.Order("id ASC")
	}

	// Pagination
	limit := filters.Limit
	if limit <= 0 {
		limit = 20
	}
	if err := tx.Limit(limit).Offset(filters.Offset).Find(&recs).Error; err != nil {
		return nil, 0, err
	}

	domainRecs := make([]resources.ResourceRecord, len(recs))
	for i, rec := range recs {
		domainRecs[i] = *toDomainRecord(&rec)
	}

	return domainRecs, int(count), nil
}

func (r *gormRepository) UpdateRecord(ctx context.Context, record *resources.ResourceRecord) (*resources.ResourceRecord, error) {
	var dbRecord models.ResourceRecord
	if err := r.db.WithContext(ctx).Where("id = ? AND collection_id = ? AND organisation_id = ?", record.ID, record.CollectionID, record.OrganisationID).First(&dbRecord).Error; err != nil {
		return nil, err
	}

	dbRecord.Data = models.RecordData(record.Data)
	dbRecord.UpdatedBy = record.UpdatedBy

	if err := r.db.WithContext(ctx).Save(&dbRecord).Error; err != nil {
		return nil, err
	}

	return toDomainRecord(&dbRecord), nil
}

func (r *gormRepository) DeleteRecord(ctx context.Context, id uint, collectionID uint, orgID uint) error {
	res := r.db.WithContext(ctx).Where("id = ? AND collection_id = ? AND organisation_id = ?", id, collectionID, orgID).Delete(&models.ResourceRecord{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *gormRepository) SearchRecords(ctx context.Context, orgID uint, query string, limit int, offset int) ([]resources.ResourceRecord, error) {
	var records []models.ResourceRecord
	err := r.db.WithContext(ctx).
		Where("organisation_id = ? AND data::text ILIKE ?", orgID, "%"+query+"%").
		Limit(limit).
		Offset(offset).
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	domainRecords := make([]resources.ResourceRecord, len(records))
	for i, rec := range records {
		domainRecords[i] = *toDomainRecord(&rec)
	}
	return domainRecords, nil
}

func sanitizeFieldName(s string) string {
	var result []rune
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result = append(result, r)
		}
	}
	return string(result)
}
