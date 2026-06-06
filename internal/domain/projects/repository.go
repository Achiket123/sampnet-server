package projects

import "context"

type Repository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uint) (*Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uint) error
	GetByOrganisation(ctx context.Context, orgID uint) ([]Project, error)
	GetByTeam(ctx context.Context, teamID uint) ([]Project, error)
	GetWithLessData(ctx context.Context, orgID uint) ([]Project, error)
	CreateMilestone(ctx context.Context, milestone *Milestone) error
	UpdateMilestone(ctx context.Context, milestone *Milestone) error
	DeleteMilestone(ctx context.Context, id uint) error
	GetMilestoneByID(ctx context.Context, id uint) (*Milestone, error)
}
