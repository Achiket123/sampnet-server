package files

import (
	"context"
	domain "server/internal/domain/files"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, file *domain.File) error {
	model := &models.File{
		FileName: file.FileName,
		Data:     file.Data,
		FileType: file.FileType,
		FileSize: file.FileSize,
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	file.ID = model.ID
	return nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*domain.File, error) {
	var model models.File
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return &domain.File{
		ID:       model.ID,
		FileName: model.FileName,
		Data:     model.Data,
		FileType: model.FileType,
		FileSize: model.FileSize,
	}, nil
}
