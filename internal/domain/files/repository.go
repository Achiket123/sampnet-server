package files

import "context"

type Repository interface {
	Create(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id uint) (*File, error)
	Upload(ctx context.Context, file *File) (string, error)
}
