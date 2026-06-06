package tasks

import (
	"context"
	"fmt"
	"log"
	"server/internal/domain/notifications"
	domain "server/internal/domain/tasks"
)

type service struct {
	repo           domain.Repository
	notificationUc notifications.UseCase
}

func NewService(repo domain.Repository, notificationUc notifications.UseCase) domain.UseCase {
	return &service{
		repo:           repo,
		notificationUc: notificationUc,
	}
}

func (s *service) CreateTask(ctx context.Context, task *domain.Task) error {
	if err := s.repo.Create(ctx, task); err != nil {
		return err
	}

	if task.AssignedTo != 0 {
		go func() {
			title := "New task assigned"
			message := fmt.Sprintf("You have been assigned a new task: %s", task.Title)
			link := fmt.Sprintf("/task-detail/%d", task.ID)
			err := s.notificationUc.CreateNotification(context.Background(), task.AssignedTo, task.OrganisationID, title, message, "task_assigned", link)
			if err != nil {
				log.Printf("failed to create notification for task assignment: %v", err)
			}
		}()
	}

	return nil
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
