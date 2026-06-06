package projects

import (
	"context"
	domain "server/internal/domain/projects"
	"time"
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
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	for i := range project.Milestones {
		project.Milestones[i].IsOverdue = project.Milestones[i].Status != "Completed" && time.Now().After(project.Milestones[i].DueDate)
	}
	return project, nil
}

func (s *service) UpdateProject(ctx context.Context, project *domain.Project) error {
	return s.repo.Update(ctx, project)
}

func (s *service) DeleteProject(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetProjectsByOrganisation(ctx context.Context, orgID uint) ([]domain.Project, error) {
	projects, err := s.repo.GetByOrganisation(ctx, orgID)
	if err != nil {
		return nil, err
	}
	for i := range projects {
		for j := range projects[i].Milestones {
			projects[i].Milestones[j].IsOverdue = projects[i].Milestones[j].Status != "Completed" && time.Now().After(projects[i].Milestones[j].DueDate)
		}
	}
	return projects, nil
}

func (s *service) GetProjectsByTeam(ctx context.Context, teamID uint) ([]domain.Project, error) {
	projects, err := s.repo.GetByTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}
	for i := range projects {
		for j := range projects[i].Milestones {
			projects[i].Milestones[j].IsOverdue = projects[i].Milestones[j].Status != "Completed" && time.Now().After(projects[i].Milestones[j].DueDate)
		}
	}
	return projects, nil
}

func (s *service) GetProjectsWithLessData(ctx context.Context, orgID uint) ([]domain.Project, error) {
	return s.repo.GetWithLessData(ctx, orgID)
}

func (s *service) CreateMilestone(ctx context.Context, milestone *domain.Milestone) error {
	return s.repo.CreateMilestone(ctx, milestone)
}

func (s *service) UpdateMilestone(ctx context.Context, milestone *domain.Milestone) error {
	return s.repo.UpdateMilestone(ctx, milestone)
}

func (s *service) DeleteMilestone(ctx context.Context, id uint) error {
	return s.repo.DeleteMilestone(ctx, id)
}

func (s *service) GetMilestoneByID(ctx context.Context, id uint) (*domain.Milestone, error) {
	milestone, err := s.repo.GetMilestoneByID(ctx, id)
	if err != nil {
		return nil, err
	}
	milestone.IsOverdue = milestone.Status != "Completed" && time.Now().After(milestone.DueDate)
	return milestone, nil
}

