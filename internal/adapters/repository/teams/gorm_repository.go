package teams

import (
	"context"
	authDomain "server/internal/domain/auth"
	orgDomain "server/internal/domain/organisation"
	domain "server/internal/domain/teams"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, team *domain.Team) error {
	model := toModel(team)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	team.ID = model.ID
	return nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*domain.Team, error) {
	var model models.Team
	if err := r.db.WithContext(ctx).Preload("Organisation").Preload("CreatedByUser").Preload("TeamLeadUser").First(&model, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) Update(ctx context.Context, team *domain.Team) error {
	model := toModel(team)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Team{}, id).Error
}

func (r *gormRepository) GetByOrganisation(ctx context.Context, orgID uint) ([]domain.Team, error) {
	var modelsList []models.Team
	if err := r.db.WithContext(ctx).Preload("Organisation").Preload("CreatedByUser").Preload("TeamLeadUser").Where("organisation_id = ?", orgID).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func (r *gormRepository) GetAll(ctx context.Context) ([]domain.Team, error) {
	var modelsList []models.Team
	if err := r.db.WithContext(ctx).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func (r *gormRepository) CreateMember(ctx context.Context, member *domain.TeamMember) error {
	model := &models.TeamMember{
		UserID:   member.UserID,
		TeamID:   member.TeamID,
		Role:     member.Role,
		IsActive: member.IsActive,
		IsLeader: member.IsLeader,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *gormRepository) GetMembersByTeam(ctx context.Context, teamID uint) ([]domain.TeamMember, error) {
	var modelsList []models.TeamMember
	if err := r.db.WithContext(ctx).Preload("User").Where("team_id = ?", teamID).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	result := make([]domain.TeamMember, len(modelsList))
	for i, m := range modelsList {
		result[i] = domain.TeamMember{
			ID:     m.ID,
			UserID: m.UserID,
			TeamID: m.TeamID,
			User: authDomain.User{
				ID:        m.User.ID,
				FirstName: m.User.FirstName,
				LastName:  m.User.LastName,
			},
			Role:     m.Role,
			IsActive: m.IsActive,
			IsLeader: m.IsLeader,
		}
	}
	return result, nil
}

func toModel(t *domain.Team) *models.Team {
	model := &models.Team{
		Name:           t.Name,
		Description:    t.Description,
		OrganisationID: t.OrganisationID,
		CreatedBy:      t.CreatedBy,
		TeamLead:       t.TeamLead,
		IsActive:       t.IsActive,
	}
	model.ID = t.ID
	return model
}

func toDomain(m *models.Team) *domain.Team {
	t := &domain.Team{
		ID:             m.ID,
		Name:           m.Name,
		Description:    m.Description,
		OrganisationID: m.OrganisationID,
		CreatedBy:      m.CreatedBy,
		TeamLead:       m.TeamLead,
		IsActive:       m.IsActive,
	}
	t.Organisation = orgDomain.Entity{ID: m.Organisation.ID, CompanyName: m.Organisation.CompanyName}
	t.TeamLeadUser = authDomain.User{ID: m.TeamLeadUser.ID, FirstName: m.TeamLeadUser.FirstName, LastName: m.TeamLeadUser.LastName}
	t.CreatedByUser = authDomain.User{ID: m.CreatedByUser.ID, FirstName: m.CreatedByUser.FirstName, LastName: m.CreatedByUser.LastName}
	return t
}

func toDomainList(modelsList []models.Team) []domain.Team {
	result := make([]domain.Team, len(modelsList))
	for i, m := range modelsList {
		result[i] = *toDomain(&m)
	}
	return result
}
