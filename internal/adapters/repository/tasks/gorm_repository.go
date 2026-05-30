package tasks

import (
	"context"
	authDomain "server/internal/domain/auth"
	domain "server/internal/domain/tasks"
	orgDomain "server/internal/domain/organisation"
	teamDomain "server/internal/domain/teams"
	projectDomain "server/internal/domain/projects"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, task *domain.Task) error {
	model := toModel(task)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	task.ID = model.ID
	return nil
}

func (r *gormRepository) Update(ctx context.Context, task *domain.Task) error {
	model := toModel(task)
	return r.db.WithContext(ctx).Model(&models.Task{}).Where("id = ?", task.ID).Updates(model).Error
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Task{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*domain.Task, error) {
	var model models.Task
	if err := r.db.WithContext(ctx).Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").Preload("Team").Preload("Project").First(&model, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) GetByTeam(ctx context.Context, userID uint) ([]domain.Task, error) {
	var modelsList []models.Task
	if err := r.db.WithContext(ctx).Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").Where("assigned_to = ?", userID).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func (r *gormRepository) GetByProject(ctx context.Context, projectID uint, page int) ([]domain.Task, error) {
	var modelsList []models.Task
	if err := r.db.WithContext(ctx).Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").
		Where("project_id = ?", projectID).Limit(20).Offset((page - 1) * 20).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func (r *gormRepository) GetPersonal(ctx context.Context, userID uint, page, pageSize int) ([]domain.Task, int64, error) {
	var modelsList []models.Task
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Task{}).Where("assigned_to = ?", userID)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").
		Limit(pageSize).Offset((page - 1) * pageSize).Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}
	return toDomainList(modelsList), count, nil
}

func (r *gormRepository) GetByOrganisation(ctx context.Context, orgID uint, page, pageSize int) ([]domain.Task, int64, error) {
	var modelsList []models.Task
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Task{}).Where("organisation_id = ?", orgID)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").
		Limit(pageSize).Offset((page - 1) * pageSize).Find(&modelsList).Error; err != nil {
		return nil, 0, err
	}
	return toDomainList(modelsList), count, nil
}

func (r *gormRepository) GetByTitle(ctx context.Context, title string) ([]domain.Task, error) {
	var modelsList []models.Task
	if err := r.db.WithContext(ctx).Preload("AssignedUser").Preload("AssignedByUser").Preload("Organisation").Where("title = ?", title).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func toModel(t *domain.Task) *models.Task {
	return &models.Task{
		Title:          t.Title,
		Description:    t.Description,
		DueDate:        t.DueDate,
		AssignedTo:     t.AssignedTo,
		AssignedBy:     t.AssignedBy,
		Type:           t.Type,
		Priority:       t.Priority,
		Status:         t.Status,
		OrganisationID: t.OrganisationID,
		IsPersonal:     t.IsPersonal,
		TeamID:         t.TeamID,
		ProjectID:      t.ProjectID,
	}
}

func toDomain(m *models.Task) *domain.Task {
	t := &domain.Task{
		ID:             m.ID,
		Title:          m.Title,
		Description:    m.Description,
		DueDate:        m.DueDate,
		AssignedTo:     m.AssignedTo,
		AssignedBy:     m.AssignedBy,
		Type:           m.Type,
		Priority:       m.Priority,
		Status:         m.Status,
		OrganisationID: m.OrganisationID,
		IsPersonal:     m.IsPersonal,
		TeamID:         m.TeamID,
		ProjectID:      m.ProjectID,
	}
	if m.AssignedUser != nil {
		t.AssignedUser = &authDomain.User{ID: m.AssignedUser.ID, FirstName: m.AssignedUser.FirstName, LastName: m.AssignedUser.LastName}
	}
	if m.AssignedByUser != nil {
		t.AssignedByUser = &authDomain.User{ID: m.AssignedByUser.ID, FirstName: m.AssignedByUser.FirstName, LastName: m.AssignedByUser.LastName}
	}
	if m.Organisation != nil {
		t.Organisation = &orgDomain.Entity{ID: m.Organisation.ID, CompanyName: m.Organisation.CompanyName}
	}
	if m.Team != nil {
		t.Team = &teamDomain.Team{ID: m.Team.ID, Name: m.Team.Name}
	}
	if m.Project != nil {
		t.Project = &projectDomain.Project{ID: m.Project.ID, Name: m.Project.Name}
	}
	return t
}

func toDomainList(modelsList []models.Task) []domain.Task {
	result := make([]domain.Task, len(modelsList))
	for i, m := range modelsList {
		result[i] = *toDomain(&m)
	}
	return result
}
