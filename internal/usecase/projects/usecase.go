package projects

import (
	"context"
	domain "server/internal/domain/projects"
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.UseCase {
	return &service{repo: repo}
}

func (s *service) CreateProject(ctx context.Context, project *domain.Project) error {
	return s.repo.Create(ctx, project)
}

func (s *service) GetProject(ctx context.Context, id uint) (*domain.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateProject(ctx context.Context, project *domain.Project) error {
	return s.repo.Update(ctx, project)
}

func (s *service) DeleteProject(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetProjectsByOrganisation(ctx context.Context, orgID uint) ([]domain.Project, error) {
	return s.repo.GetByOrganisation(ctx, orgID)
}

func (s *service) GetProjectsByTeam(ctx context.Context, teamID uint) ([]domain.Project, error) {
	return s.repo.GetByTeam(ctx, teamID)
}

func (s *service) GetProjectsWithLessData(ctx context.Context, orgID uint) ([]domain.Project, error) {
	return s.repo.GetWithLessData(ctx, orgID)
}
