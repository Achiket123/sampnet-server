package files

import (
	"bytes"
	"context"
	domain "server/internal/domain/files"
	"server/internal/platform/database/models"

	cld "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"gorm.io/gorm"
)

type gormRepository struct {
	db         *gorm.DB
	cloudinary *cld.Cloudinary
}

func NewGormRepository(db *gorm.DB, cloudinaryClient *cld.Cloudinary) domain.Repository {
	return &gormRepository{
		db:         db,
		cloudinary: cloudinaryClient,
	}
}

func (r *gormRepository) Upload(ctx context.Context, file *domain.File) (string, error) {
	uploadParams := uploader.UploadParams{
		PublicID: file.FileName,
	}
	uploadResult, err := r.cloudinary.Upload.Upload(ctx, bytes.NewReader(file.Data), uploadParams)
	if err != nil {
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func (r *gormRepository) Create(ctx context.Context, file *domain.File) error {
	model := &models.File{
		FileName: file.FileName,
		URL:      file.URL,
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
		URL:      model.URL,
		FileType: model.FileType,
		FileSize: model.FileSize,
	}, nil
}
