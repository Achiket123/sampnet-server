package projects

import "context"

type UseCase interface {
	CreateProject(ctx context.Context, project *Project) error
	GetProject(ctx context.Context, id uint) (*Project, error)
	UpdateProject(ctx context.Context, project *Project) error
	DeleteProject(ctx context.Context, id uint) error
	GetProjectsByOrganisation(ctx context.Context, orgID uint) ([]Project, error)
	GetProjectsByTeam(ctx context.Context, teamID uint) ([]Project, error)
	GetProjectsWithLessData(ctx context.Context, orgID uint) ([]Project, error)
}
