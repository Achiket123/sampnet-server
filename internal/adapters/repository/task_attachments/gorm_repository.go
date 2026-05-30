package task_attachments

import (
	"context"
	"server/internal/domain/task_attachments"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) task_attachments.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, attachment *task_attachments.TaskAttachment) error {
	model := toModel(attachment)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	
	// Refetch to populate user info
	if err := r.db.WithContext(ctx).Preload("UploadedByUser").First(model, model.ID).Error; err != nil {
		return err
	}
	
	*attachment = *toDomain(model)
	return nil
}

func (r *gormRepository) GetByTaskID(ctx context.Context, taskID uint) ([]task_attachments.TaskAttachment, error) {
	var modelsList []models.TaskAttachment
	if err := r.db.WithContext(ctx).Preload("File").Preload("UploadedByUser").Where("task_id = ?", taskID).Order("created_at desc").Find(&modelsList).Error; err != nil {
		return nil, err
	}
	
	result := make([]task_attachments.TaskAttachment, len(modelsList))
	for i, m := range modelsList {
		result[i] = *toDomain(&m)
	}
	return result, nil
}

func (r *gormRepository) GetByID(ctx context.Context, attachmentID uint) (*task_attachments.TaskAttachment, error) {
	var model models.TaskAttachment
	if err := r.db.WithContext(ctx).Preload("File").Preload("UploadedByUser").First(&model, attachmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) Delete(ctx context.Context, attachmentID uint) error {
	return r.db.WithContext(ctx).Delete(&models.TaskAttachment{}, attachmentID).Error
}

func toModel(a *task_attachments.TaskAttachment) *models.TaskAttachment {
	return &models.TaskAttachment{
		ID:         a.ID,
		TaskID:     a.TaskID,
		FileID:     a.FileID,
		UploadedBy: a.UploadedBy,
		FileName:   a.FileName,
		CreatedAt:  a.CreatedAt,
	}
}

func toDomain(m *models.TaskAttachment) *task_attachments.TaskAttachment {
	return &task_attachments.TaskAttachment{
		ID:                m.ID,
		TaskID:            m.TaskID,
		FileID:            m.FileID,
		UploadedBy:        m.UploadedBy,
		FileName:          m.FileName,
		CreatedAt:         m.CreatedAt,
		UploaderFirstName: m.UploadedByUser.FirstName,
		UploaderLastName:  m.UploadedByUser.LastName,
	}
}
