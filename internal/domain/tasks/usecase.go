package tasks

import "context"

type UseCase interface {
	CreateTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, task *Task) (*Task, error)
	DeleteTask(ctx context.Context, id uint) error
	GetTask(ctx context.Context, id uint) (*Task, error)
	GetTeamTasks(ctx context.Context, userID uint) ([]Task, error)
	GetProjectTasks(ctx context.Context, projectID uint, page int) ([]Task, error)
	GetPersonalTasks(ctx context.Context, userID uint, page, pageSize int) ([]Task, int64, error)
	GetOrganisationTasks(ctx context.Context, orgID uint, page, pageSize int) ([]Task, int64, error)
	GetTaskByTitle(ctx context.Context, title string) ([]Task, error)
}
