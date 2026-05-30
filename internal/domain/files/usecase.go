package files

import "context"

type UseCase interface {
	UploadFile(ctx context.Context, file *File) error
	GetFile(ctx context.Context, id uint) (*File, error)
}
