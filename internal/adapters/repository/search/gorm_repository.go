package search

import (
	"context"
	"fmt"
	"strings"

	domain "server/internal/domain/search"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a search Repository backed by PostgreSQL via GORM.
func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

// contains checks whether the given slice contains the target string.
func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// wantType reports whether the given type should be queried given the Types filter.
// An empty types slice means "query all".
func wantType(filters *domain.SearchFilters, t string) bool {
	return len(filters.Types) == 0 || contains(filters.Types, t)
}

func (r *gormRepository) Search(ctx context.Context, filters *domain.SearchFilters) (*domain.SearchResults, error) {
	like := "%" + filters.Query + "%"
	orgID := filters.OrganisationID

	var taskItems, projectItems, teamItems, employeeItems, chatItems []domain.SearchResultItem

	// ── Tasks ────────────────────────────────────────────────────────────────
	if wantType(filters, "task") {
		var tasks []models.Task
		err := r.db.WithContext(ctx).
			Select("id, title, description, status, priority, project_id, organisation_id").
			Where("organisation_id = ? AND (title ILIKE ? OR description ILIKE ?)", orgID, like, like).
			Find(&tasks).Error
		if err != nil {
			return nil, fmt.Errorf("search tasks: %w", err)
		}
		for _, t := range tasks {
			extra := map[string]interface{}{}
			if t.ProjectID != nil {
				extra["project_id"] = *t.ProjectID
			}
			taskItems = append(taskItems, domain.SearchResultItem{
				ID:             t.ID,
				Type:           "task",
				Title:          t.Title,
				Subtitle:       t.Status + " · " + t.Priority,
				OrganisationID: t.OrganisationID,
				ExtraData:      extra,
			})
		}
	}

	// ── Projects ─────────────────────────────────────────────────────────────
	if wantType(filters, "project") {
		var projects []models.Project
		err := r.db.WithContext(ctx).
			Select("id, name, description, status, organisation_id").
			Where("organisation_id = ? AND (name ILIKE ? OR description ILIKE ?)", orgID, like, like).
			Find(&projects).Error
		if err != nil {
			return nil, fmt.Errorf("search projects: %w", err)
		}
		for _, p := range projects {
			projectItems = append(projectItems, domain.SearchResultItem{
				ID:             p.ID,
				Type:           "project",
				Title:          p.Name,
				Subtitle:       p.Status,
				OrganisationID: p.OrganisationID,
			})
		}
	}

	// ── Teams ────────────────────────────────────────────────────────────────
	if wantType(filters, "team") {
		var teams []models.Team
		err := r.db.WithContext(ctx).
			Select("id, name, description, organisation_id").
			Where("organisation_id = ? AND (name ILIKE ? OR description ILIKE ?)", orgID, like, like).
			Find(&teams).Error
		if err != nil {
			return nil, fmt.Errorf("search teams: %w", err)
		}
		for _, t := range teams {
			subtitle := t.Description
			if len(subtitle) > 80 {
				subtitle = subtitle[:80] + "…"
			}
			teamItems = append(teamItems, domain.SearchResultItem{
				ID:             t.ID,
				Type:           "team",
				Title:          t.Name,
				Subtitle:       subtitle,
				OrganisationID: t.OrganisationID,
			})
		}
	}

	// ── Employees (join employees → users) ───────────────────────────────────
	if wantType(filters, "employee") {
		type empRow struct {
			UserID    uint
			FirstName string
			LastName  string
			Email     string
			Type      string
		}
		var rows []empRow
		err := r.db.WithContext(ctx).
			Table("employees").
			Select("employees.user_id, user_models.first_name, user_models.last_name, user_models.email, employees.type").
			Joins("JOIN user_models ON user_models.id = employees.user_id").
			Where("employees.organisation_id = ? AND (user_models.first_name ILIKE ? OR user_models.last_name ILIKE ? OR user_models.email ILIKE ?)", orgID, like, like, like).
			Scan(&rows).Error
		if err != nil {
			return nil, fmt.Errorf("search employees: %w", err)
		}
		for _, e := range rows {
			employeeItems = append(employeeItems, domain.SearchResultItem{
				ID:             e.UserID,
				Type:           "employee",
				Title:          strings.TrimSpace(e.FirstName + " " + e.LastName),
				Subtitle:       e.Type,
				OrganisationID: orgID,
				ExtraData: map[string]interface{}{
					"user_id": e.UserID,
					"email":   e.Email,
				},
			})
		}
	}

	// ── Chats (group chats only — Chat model stores DMs by user name) ────────
	// NOTE: The Chat model in this codebase is a DM/sync table (not a group chat
	// model with is_group). We skip the is_group filter and search by first_name
	// or last_name within the chat table so managers can find colleagues quickly.
	// If a true group-chat model is introduced later, add is_group = true here.
	if wantType(filters, "chat") {
		var chats []models.Chat
		err := r.db.WithContext(ctx).
			Select("id, first_name, last_name, organisation_id").
			Where("organisation_id = ? AND (first_name ILIKE ? OR last_name ILIKE ?)", orgID, like, like).
			Find(&chats).Error
		if err != nil {
			return nil, fmt.Errorf("search chats: %w", err)
		}
		for _, c := range chats {
			chatItems = append(chatItems, domain.SearchResultItem{
				ID:             c.ID,
				Type:           "chat",
				Title:          strings.TrimSpace(c.FirstName + " " + c.LastName),
				OrganisationID: c.OrganisationID,
			})
		}
	}

	// ── Resources ────────────────────────────────────────────────────────────
	var resourceItems []domain.SearchResultItem
	if wantType(filters, "resource") {
		type resourceRow struct {
			ID             uint
			CollectionID   uint
			CollectionName string
			Data           models.RecordData
			OrganisationID uint
		}
		var rows []resourceRow
		err := r.db.WithContext(ctx).
			Table("resource_records").
			Select("resource_records.id, resource_records.collection_id, resource_collections.name as collection_name, resource_records.data, resource_records.organisation_id").
			Joins("JOIN resource_collections ON resource_collections.id = resource_records.collection_id").
			Where("resource_records.organisation_id = ? AND resource_records.data::text ILIKE ? AND resource_records.deleted_at IS NULL", orgID, like).
			Scan(&rows).Error
		if err != nil {
			return nil, fmt.Errorf("search resources: %w", err)
		}
		for _, row := range rows {
			title := ""
			for _, v := range row.Data {
				if strVal, ok := v.(string); ok && strVal != "" {
					title = strVal
					break
				}
			}
			if title == "" {
				title = fmt.Sprintf("Record #%d", row.ID)
			}
			resourceItems = append(resourceItems, domain.SearchResultItem{
				ID:             row.ID,
				Type:           "resource",
				Title:          title,
				Subtitle:       row.CollectionName,
				OrganisationID: row.OrganisationID,
				ExtraData: map[string]interface{}{
					"collection_id": row.CollectionID,
				},
			})
		}
	}
	
	// ── Research ─────────────────────────────────────────────────────────────
	var researchItems []domain.SearchResultItem
	if wantType(filters, "research") {
		var researches []models.ResearchEntry
		err := r.db.WithContext(ctx).
			Select("id, title, description, status, organisation_id").
			Where("organisation_id = ? AND (title ILIKE ? OR description ILIKE ?) AND deleted_at IS NULL", orgID, like, like).
			Find(&researches).Error
		if err != nil {
			return nil, fmt.Errorf("search researches: %w", err)
		}
		for _, re := range researches {
			researchItems = append(researchItems, domain.SearchResultItem{
				ID:             re.ID,
				Type:           "research",
				Title:          re.Title,
				Subtitle:       re.Status,
				OrganisationID: re.OrganisationID,
			})
		}
	}

	// ── Merge all results and apply global pagination ─────────────────────────
	all := make([]domain.SearchResultItem, 0,
		len(taskItems)+len(projectItems)+len(teamItems)+len(employeeItems)+len(chatItems)+len(resourceItems)+len(researchItems))
	all = append(all, taskItems...)
	all = append(all, projectItems...)
	all = append(all, teamItems...)
	all = append(all, employeeItems...)
	all = append(all, chatItems...)
	all = append(all, resourceItems...)
	all = append(all, researchItems...)

	// Apply offset
	if filters.Offset >= len(all) {
		all = []domain.SearchResultItem{}
	} else {
		all = all[filters.Offset:]
	}
	// Apply limit
	if filters.Limit > 0 && len(all) > filters.Limit {
		all = all[:filters.Limit]
	}

	return &domain.SearchResults{
		Query:         filters.Query,
		TotalCount:    len(all),
		Items:         all,
		TaskCount:     len(taskItems),
		ProjectCount:  len(projectItems),
		TeamCount:     len(teamItems),
		EmployeeCount: len(employeeItems),
		ChatCount:     len(chatItems),
		ResourceCount: len(resourceItems),
		ResearchCount: len(researchItems),
	}, nil
}