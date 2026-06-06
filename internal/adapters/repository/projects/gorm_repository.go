package projects

import (
	"context"
	authDomain "server/internal/domain/auth"
	domain "server/internal/domain/projects"
	orgDomain "server/internal/domain/organisation"
	teamDomain "server/internal/domain/teams"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, project *domain.Project) error {
	model := toModel(project)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	project.ID = model.ID
	return nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*domain.Project, error) {
	var model models.Project
	if err := r.db.WithContext(ctx).Preload("Team").Preload("CreatedByUser").Preload("Organisation").Preload("Milestones").First(&model, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) Update(ctx context.Context, project *domain.Project) error {
	model := toModel(project)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, id).Error
}

func (r *gormRepository) GetByOrganisation(ctx context.Context, orgID uint) ([]domain.Project, error) {
	var modelsList []models.Project
	if err := r.db.WithContext(ctx).Preload("Team").Preload("CreatedByUser").Preload("Milestones").Where("organisation_id = ?", orgID).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func (r *gormRepository) GetByTeam(ctx context.Context, teamID uint) ([]domain.Project, error) {
	var modelsList []models.Project
	if err := r.db.WithContext(ctx).Preload("Team").Preload("CreatedByUser").Preload("Milestones").Where("team_id = ?", teamID).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func (r *gormRepository) GetWithLessData(ctx context.Context, orgID uint) ([]domain.Project, error) {
	var modelsList []models.Project
	if err := r.db.WithContext(ctx).Where("organisation_id = ?", orgID).Select("id", "name", "description").Find(&modelsList).Error; err != nil {
		return nil, err
	}
	result := make([]domain.Project, len(modelsList))
	for i, m := range modelsList {
		result[i] = domain.Project{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
		}
	}
	return result, nil
}

func (r *gormRepository) CreateMilestone(ctx context.Context, milestone *domain.Milestone) error {
	model := &models.Milestone{
		ProjectID:      milestone.ProjectID,
		Title:          milestone.Title,
		Description:    milestone.Description,
		DueDate:        milestone.DueDate,
		Status:         milestone.Status,
		OrganisationID: milestone.OrganisationID,
	}
	model.ID = milestone.ID
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	milestone.ID = model.ID
	return nil
}

func (r *gormRepository) UpdateMilestone(ctx context.Context, milestone *domain.Milestone) error {
	model := &models.Milestone{
		ProjectID:      milestone.ProjectID,
		Title:          milestone.Title,
		Description:    milestone.Description,
		DueDate:        milestone.DueDate,
		Status:         milestone.Status,
		OrganisationID: milestone.OrganisationID,
	}
	model.ID = milestone.ID
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *gormRepository) DeleteMilestone(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Milestone{}, id).Error
}

func (r *gormRepository) GetMilestoneByID(ctx context.Context, id uint) (*domain.Milestone, error) {
	var model models.Milestone
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return &domain.Milestone{
		ID:             model.ID,
		ProjectID:      model.ProjectID,
		Title:          model.Title,
		Description:    model.Description,
		DueDate:        model.DueDate,
		Status:         model.Status,
		OrganisationID: model.OrganisationID,
	}, nil
}

func toModel(p *domain.Project) *models.Project {
	model := &models.Project{
		Name:             p.Name,
		Description:      p.Description,
		StartDate:        p.StartDate,
		EndDate:          p.EndDate,
		OrganisationID:   p.OrganisationID,
		TeamID:           p.TeamID,
		CreatedBy:        p.CreatedBy,
		Status:           p.Status,
		Priority:         p.Priority,
		CompletionStatus: p.CompletionStatus,
	}
	model.ID = p.ID
	return model
}

func toDomain(m *models.Project) *domain.Project {
	p := &domain.Project{
		ID:               m.ID,
		Name:             m.Name,
		Description:      m.Description,
		StartDate:        m.StartDate,
		EndDate:          m.EndDate,
		OrganisationID:   m.OrganisationID,
		TeamID:           m.TeamID,
		CreatedBy:        m.CreatedBy,
		Status:           m.Status,
		Priority:         m.Priority,
		CompletionStatus: m.CompletionStatus,
	}
	p.Organisation = orgDomain.Entity{ID: m.Organisation.ID, CompanyName: m.Organisation.CompanyName}
	p.Team = teamDomain.Team{ID: m.Team.ID, Name: m.Team.Name}
	p.CreatedByUser = authDomain.User{ID: m.CreatedByUser.ID, FirstName: m.CreatedByUser.FirstName, LastName: m.CreatedByUser.LastName}
	
	if len(m.Milestones) > 0 {
		p.Milestones = make([]domain.Milestone, len(m.Milestones))
		for i, ms := range m.Milestones {
			p.Milestones[i] = domain.Milestone{
				ID:             ms.ID,
				ProjectID:      ms.ProjectID,
				Title:          ms.Title,
				Description:    ms.Description,
				DueDate:        ms.DueDate,
				Status:         ms.Status,
				OrganisationID: ms.OrganisationID,
			}
		}
	}
	return p
}

func toDomainList(modelsList []models.Project) []domain.Project {
	result := make([]domain.Project, len(modelsList))
	for i, m := range modelsList {
		result[i] = *toDomain(&m)
	}
	return result
}

