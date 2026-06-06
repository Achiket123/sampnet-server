package resources

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"server/internal/domain/resources"
)

var (
	keyRegex   = regexp.MustCompile(`^[a-z0-9_]+$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type service struct {
	repo resources.Repository
}

func NewService(repo resources.Repository) resources.UseCase {
	return &service{repo: repo}
}

func (s *service) CreateCollection(ctx context.Context, orgID uint, createdBy uint, name string, description string, icon *string, colour *string, fields []resources.FieldDefinition) (*resources.ResourceCollection, error) {
	if name == "" {
		return nil, &resources.ValidationError{Errors: map[string]string{"name": "Collection name is required"}}
	}
	if len(name) > 100 {
		return nil, &resources.ValidationError{Errors: map[string]string{"name": "Collection name cannot exceed 100 characters"}}
	}

	// Check name uniqueness in the org
	existingColls, _, err := s.repo.GetCollectionsByOrg(ctx, orgID, 1000, 0)
	if err == nil {
		for _, coll := range existingColls {
			if strings.EqualFold(coll.Name, name) {
				return nil, &resources.ValidationError{Errors: map[string]string{"name": "A collection with this name already exists in your organisation"}}
			}
		}
	}

	// Validate fields
	validationErrs := make(map[string]string)
	for i, f := range fields {
		if f.Key == "" || !keyRegex.MatchString(f.Key) {
			validationErrs[fmt.Sprintf("fields[%d].key", i)] = "Field key must be snake_case (alphanumeric and underscores only)"
		}
		if f.Label == "" {
			validationErrs[fmt.Sprintf("fields[%d].label", i)] = "Field label is required"
		}
		switch f.Type {
		case "string", "text", "number", "boolean", "date", "url", "email", "phone", "select", "multi_select":

		default:
			validationErrs[fmt.Sprintf("fields[%d].type", i)] = "Invalid field type"
		}
		if (f.Type == "select" || f.Type == "multi_select") && len(f.Options) == 0 {
			validationErrs[fmt.Sprintf("fields[%d].options", i)] = "Options are required for select and multi_select fields"
		}
	}

	if len(validationErrs) > 0 {
		return nil, &resources.ValidationError{Errors: validationErrs}
	}

	collection := &resources.ResourceCollection{
		OrganisationID: orgID,
		Name:           name,
		Description:    description,
		Icon:           icon,
		Colour:         colour,
		CreatedBy:      createdBy,
		Fields:         fields,
	}

	return s.repo.CreateCollection(ctx, collection)
}

func (s *service) GetCollection(ctx context.Context, id uint, orgID uint) (*resources.ResourceCollection, error) {
	return s.repo.GetCollectionByID(ctx, id, orgID)
}

func (s *service) ListCollections(ctx context.Context, orgID uint, limit int, offset int) ([]resources.ResourceCollection, int, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetCollectionsByOrg(ctx, orgID, limit, offset)
}

func (s *service) UpdateCollection(ctx context.Context, id uint, orgID uint, name string, description string, icon *string, colour *string) (*resources.ResourceCollection, error) {
	if name == "" {
		return nil, &resources.ValidationError{Errors: map[string]string{"name": "Collection name is required"}}
	}

	existingColls, _, err := s.repo.GetCollectionsByOrg(ctx, orgID, 1000, 0)
	if err == nil {
		for _, coll := range existingColls {
			if coll.ID != id && strings.EqualFold(coll.Name, name) {
				return nil, &resources.ValidationError{Errors: map[string]string{"name": "A collection with this name already exists in your organisation"}}
			}
		}
	}

	collection := &resources.ResourceCollection{
		ID:             id,
		OrganisationID: orgID,
		Name:           name,
		Description:    description,
		Icon:           icon,
		Colour:         colour,
	}

	return s.repo.UpdateCollection(ctx, collection)
}

func (s *service) DeleteCollection(ctx context.Context, id uint, orgID uint, force bool) error {
	if !force {
		// Check if records exist
		recs, count, err := s.repo.GetRecordsByCollection(ctx, id, orgID, resources.RecordFilters{Limit: 1})
		if err == nil && count > 0 {
			return fmt.Errorf("collection has %d record(s); requires force=true query param to delete", count)
		}
		_ = recs
	}
	return s.repo.DeleteCollection(ctx, id, orgID)
}

func (s *service) AddFieldToCollection(ctx context.Context, collectionID uint, orgID uint, field resources.FieldDefinition) (*resources.ResourceCollection, error) {
	if field.Key == "" || !keyRegex.MatchString(field.Key) {
		return nil, &resources.ValidationError{Errors: map[string]string{"key": "Field key must be snake_case (alphanumeric and underscores)"}}
	}
	if field.Label == "" {
		return nil, &resources.ValidationError{Errors: map[string]string{"label": "Field label is required"}}
	}
	switch field.Type {
	case "text", "number", "boolean", "date", "url", "email", "phone", "select", "multi_select":
		// valid
	default:
		return nil, &resources.ValidationError{Errors: map[string]string{"type": "Invalid field type"}}
	}
	if (field.Type == "select" || field.Type == "multi_select") && len(field.Options) == 0 {
		return nil, &resources.ValidationError{Errors: map[string]string{"options": "Options are required for select and multi_select types"}}
	}

	return s.repo.AddField(ctx, collectionID, orgID, field)
}

func (s *service) UpdateFieldInCollection(ctx context.Context, collectionID uint, orgID uint, fieldKey string, field resources.FieldDefinition) (*resources.ResourceCollection, error) {
	coll, err := s.repo.GetCollectionByID(ctx, collectionID, orgID)
	if err != nil {
		return nil, err
	}

	// Find the existing field
	var targetField *resources.FieldDefinition
	for _, f := range coll.Fields {
		if f.Key == fieldKey {
			targetField = &f
			break
		}
	}

	if targetField == nil {
		return nil, errors.New("field key not found")
	}

	// Enforce that field type cannot change
	if field.Type != "" && field.Type != targetField.Type {
		return nil, errors.New("cannot change field type of an existing field")
	}

	if field.Label == "" {
		return nil, &resources.ValidationError{Errors: map[string]string{"label": "Field label is required"}}
	}
	if (targetField.Type == "select" || targetField.Type == "multi_select") && len(field.Options) == 0 {
		return nil, &resources.ValidationError{Errors: map[string]string{"options": "Options cannot be empty for select or multi_select fields"}}
	}

	return s.repo.UpdateField(ctx, collectionID, orgID, fieldKey, field)
}

func (s *service) RemoveFieldFromCollection(ctx context.Context, collectionID uint, orgID uint, fieldKey string) (*resources.ResourceCollection, string, error) {
	// Check if records have data for this field
	recs, count, err := s.repo.GetRecordsByCollection(ctx, collectionID, orgID, resources.RecordFilters{Limit: 1000})
	hasDataCount := 0
	if err == nil && count > 0 {
		for _, rec := range recs {
			if val, exists := rec.Data[fieldKey]; exists && val != nil && val != "" {
				hasDataCount++
			}
		}
	}

	updatedColl, err := s.repo.RemoveField(ctx, collectionID, orgID, fieldKey)
	if err != nil {
		return nil, "", err
	}

	var warning string
	if hasDataCount > 0 {
		warning = fmt.Sprintf("%d records have data for this field. Removing it hides it from the UI but does not delete historical data.", hasDataCount)
	}

	return updatedColl, warning, nil
}

func (s *service) validateRecordData(fields []resources.FieldDefinition, data map[string]interface{}) map[string]string {
	errs := make(map[string]string)

	for _, f := range fields {
		val, exists := data[f.Key]

		// Check required field
		if f.Required {
			if !exists || val == nil || val == "" {
				errs[f.Key] = fmt.Sprintf("Field '%s' is required", f.Label)
				continue
			}
		}

		if !exists || val == nil || val == "" {
			continue
		}

		// Type validation
		switch f.Type {
		case "number":
			switch val.(type) {
			case float64, float32, int, int64, int32:
				// ok
			case string:
				if _, err := strconv.ParseFloat(val.(string), 64); err != nil {
					errs[f.Key] = "Must be a valid number"
				}
			default:
				errs[f.Key] = "Must be a number"
			}
		case "boolean":
			switch val.(type) {
			case bool:
				// ok
			case string:
				if val.(string) != "true" && val.(string) != "false" {
					errs[f.Key] = "Must be a boolean value"
				}
			default:
				errs[f.Key] = "Must be a boolean"
			}
		case "date":
			if str, ok := val.(string); ok {
				if _, err := time.Parse(time.RFC3339, str); err != nil {
					if _, err2 := time.Parse("2006-01-02", str); err2 != nil {
						errs[f.Key] = "Must be a valid RFC3339 date string or YYYY-MM-DD"
					}
				}
			} else {
				errs[f.Key] = "Must be a string representing a date"
			}
		case "email":
			if str, ok := val.(string); ok {
				if !emailRegex.MatchString(str) {
					errs[f.Key] = "Must be a valid email address"
				}
			} else {
				errs[f.Key] = "Must be a string"
			}
		case "select":
			if str, ok := val.(string); ok {
				found := false
				for _, opt := range f.Options {
					if opt == str {
						found = true
						break
					}
				}
				if !found {
					errs[f.Key] = fmt.Sprintf("Must be one of the options: %v", f.Options)
				}
			} else {
				errs[f.Key] = "Must be a string option"
			}
		case "multi_select":
			if rawSlice, ok := val.([]interface{}); ok {
				for _, item := range rawSlice {
					strItem, isStr := item.(string)
					if !isStr {
						errs[f.Key] = "All multi-select items must be strings"
						break
					}
					found := false
					for _, opt := range f.Options {
						if opt == strItem {
							found = true
							break
						}
					}
					if !found {
						errs[f.Key] = fmt.Sprintf("Item '%s' is not a valid option", strItem)
						break
					}
				}
			} else {
				errs[f.Key] = "Must be an array of selected options"
			}
		}
	}

	return errs
}

func (s *service) CreateRecord(ctx context.Context, orgID uint, collectionID uint, createdBy uint, data map[string]interface{}, projectID *uint, teamID *uint, taskID *uint) (*resources.ResourceRecord, error) {
	coll, err := s.repo.GetCollectionByID(ctx, collectionID, orgID)
	if err != nil {
		return nil, err
	}

	// Validate against schema
	validationErrs := s.validateRecordData(coll.Fields, data)
	if len(validationErrs) > 0 {
		return nil, &resources.ValidationError{Errors: validationErrs}
	}

	// Apply default values
	for _, f := range coll.Fields {
		if _, exists := data[f.Key]; !exists || data[f.Key] == nil {
			if f.DefaultValue != nil {
				data[f.Key] = f.DefaultValue
			}
		}
	}

	record := &resources.ResourceRecord{
		CollectionID:   collectionID,
		OrganisationID: orgID,
		Data:           data,
		CreatedBy:      createdBy,
		ProjectID:      projectID,
		TeamID:         teamID,
		TaskID:         taskID,
	}

	return s.repo.CreateRecord(ctx, record)
}

func (s *service) GetRecord(ctx context.Context, id uint, collectionID uint, orgID uint) (*resources.ResourceRecord, error) {
	return s.repo.GetRecordByID(ctx, id, collectionID, orgID)
}

func (s *service) ListRecords(ctx context.Context, collectionID uint, orgID uint, filters resources.RecordFilters) ([]resources.ResourceRecord, int, error) {
	return s.repo.GetRecordsByCollection(ctx, collectionID, orgID, filters)
}

func (s *service) UpdateRecord(ctx context.Context, id uint, collectionID uint, orgID uint, updatedBy uint, data map[string]interface{}) (*resources.ResourceRecord, error) {
	coll, err := s.repo.GetCollectionByID(ctx, collectionID, orgID)
	if err != nil {
		return nil, err
	}

	existingRecord, err := s.repo.GetRecordByID(ctx, id, collectionID, orgID)
	if err != nil {
		return nil, err
	}

	// Merge patch: copy over updated keys
	mergedData := make(map[string]interface{})
	for k, v := range existingRecord.Data {
		mergedData[k] = v
	}
	for k, v := range data {
		mergedData[k] = v
	}

	// Validate merged data
	validationErrs := s.validateRecordData(coll.Fields, mergedData)
	if len(validationErrs) > 0 {
		return nil, &resources.ValidationError{Errors: validationErrs}
	}

	existingRecord.Data = mergedData
	existingRecord.UpdatedBy = &updatedBy

	return s.repo.UpdateRecord(ctx, existingRecord)
}

func (s *service) DeleteRecord(ctx context.Context, id uint, collectionID uint, orgID uint) error {
	return s.repo.DeleteRecord(ctx, id, collectionID, orgID)
}

func (s *service) BulkCreateRecords(ctx context.Context, orgID uint, collectionID uint, createdBy uint, records []map[string]interface{}) (*resources.BulkCreateResult, error) {
	coll, err := s.repo.GetCollectionByID(ctx, collectionID, orgID)
	if err != nil {
		return nil, err
	}

	result := &resources.BulkCreateResult{
		Errors: make([]resources.BulkCreateRecordError, 0),
	}

	for i, rdata := range records {
		validationErrs := s.validateRecordData(coll.Fields, rdata)
		if len(validationErrs) > 0 {
			result.Failed++
			for k, v := range validationErrs {
				result.Errors = append(result.Errors, resources.BulkCreateRecordError{
					Row:     i + 1,
					Field:   k,
					Message: v,
				})
			}
			continue
		}

		// Apply defaults
		for _, f := range coll.Fields {
			if _, exists := rdata[f.Key]; !exists || rdata[f.Key] == nil {
				if f.DefaultValue != nil {
					rdata[f.Key] = f.DefaultValue
				}
			}
		}

		record := &resources.ResourceRecord{
			CollectionID:   collectionID,
			OrganisationID: orgID,
			Data:           rdata,
			CreatedBy:      createdBy,
		}

		_, err := s.repo.CreateRecord(ctx, record)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, resources.BulkCreateRecordError{
				Row:     i + 1,
				Field:   "_db",
				Message: err.Error(),
			})
		} else {
			result.Created++
		}
	}

	return result, nil
}

func (s *service) ExportRecords(ctx context.Context, collectionID uint, orgID uint) ([]map[string]interface{}, []string, error) {
	coll, err := s.repo.GetCollectionByID(ctx, collectionID, orgID)
	if err != nil {
		return nil, nil, err
	}

	recs, _, err := s.repo.GetRecordsByCollection(ctx, collectionID, orgID, resources.RecordFilters{Limit: 100000})
	if err != nil {
		return nil, nil, err
	}

	headers := make([]string, len(coll.Fields))
	for i, f := range coll.Fields {
		headers[i] = f.Key
	}

	exported := make([]map[string]interface{}, len(recs))
	for i, rec := range recs {
		flat := make(map[string]interface{})
		for _, f := range coll.Fields {
			if val, exists := rec.Data[f.Key]; exists {
				flat[f.Key] = val
			} else {
				flat[f.Key] = ""
			}
		}
		exported[i] = flat
	}

	return exported, headers, nil
}
