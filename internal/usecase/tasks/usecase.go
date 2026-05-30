package tasks

import (
	"context"
	domain "server/internal/domain/tasks"
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.UseCase {
	return &service{repo: repo}
}

func (s *service) CreateTask(ctx context.Context, task *domain.Task) error {
	return s.repo.Create(ctx, task)
}

func (s *service) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, task.ID)
}

func (s *service) DeleteTask(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetTask(ctx context.Context, id uint) (*domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetTeamTasks(ctx context.Context, userID uint) ([]domain.Task, error) {
	return s.repo.GetByTeam(ctx, userID)
}

func (s *service) GetProjectTasks(ctx context.Context, projectID uint, page int) ([]domain.Task, error) {
	return s.repo.GetByProject(ctx, projectID, page)
}

func (s *service) GetPersonalTasks(ctx context.Context, userID uint, page, pageSize int) ([]domain.Task, int64, error) {
	return s.repo.GetPersonal(ctx, userID, page, pageSize)
}

func (s *service) GetOrganisationTasks(ctx context.Context, orgID uint, page, pageSize int) ([]domain.Task, int64, error) {
	return s.repo.GetByOrganisation(ctx, orgID, page, pageSize)
}

func (s *service) GetTaskByTitle(ctx context.Context, title string) ([]domain.Task, error) {
	return s.repo.GetByTitle(ctx, title)
}
