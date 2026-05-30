package files

import (
	"context"
	domain "server/internal/domain/files"
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.UseCase {
	return &service{repo: repo}
}

func (s *service) UploadFile(ctx context.Context, file *domain.File) error {
	return s.repo.Create(ctx, file)
}

func (s *service) GetFile(ctx context.Context, id uint) (*domain.File, error) {
	return s.repo.GetByID(ctx, id)
}
