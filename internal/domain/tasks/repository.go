package tasks

import "context"

type Repository interface {
	Create(ctx context.Context, task *Task) error
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*Task, error)
	GetByTeam(ctx context.Context, userID uint) ([]Task, error)
	GetByProject(ctx context.Context, projectID uint, page int) ([]Task, error)
	GetPersonal(ctx context.Context, userID uint, page, pageSize int) ([]Task, int64, error)
	GetByOrganisation(ctx context.Context, orgID uint, page, pageSize int) ([]Task, int64, error)
	GetByTitle(ctx context.Context, title string) ([]Task, error)
}
