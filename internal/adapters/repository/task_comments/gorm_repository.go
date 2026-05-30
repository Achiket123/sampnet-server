package task_comments

import (
	"context"
	"server/internal/domain/task_comments"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) task_comments.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, comment *task_comments.TaskComment) error {
	model := toModel(comment)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	
	// Refetch to populate user info
	if err := r.db.WithContext(ctx).Preload("User").First(model, model.ID).Error; err != nil {
		return err
	}
	
	*comment = *toDomain(model)
	return nil
}

func (r *gormRepository) GetByTaskID(ctx context.Context, taskID uint) ([]task_comments.TaskComment, error) {
	var modelsList []models.TaskComment
	if err := r.db.WithContext(ctx).Preload("User").Where("task_id = ?", taskID).Order("created_at asc").Find(&modelsList).Error; err != nil {
		return nil, err
	}
	
	result := make([]task_comments.TaskComment, len(modelsList))
	for i, m := range modelsList {
		result[i] = *toDomain(&m)
	}
	return result, nil
}

func (r *gormRepository) GetByID(ctx context.Context, commentID uint) (*task_comments.TaskComment, error) {
	var model models.TaskComment
	if err := r.db.WithContext(ctx).Preload("User").First(&model, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) Delete(ctx context.Context, commentID uint) error {
	return r.db.WithContext(ctx).Delete(&models.TaskComment{}, commentID).Error
}

func toModel(c *task_comments.TaskComment) *models.TaskComment {
	return &models.TaskComment{
		ID:        c.ID,
		TaskID:    c.TaskID,
		UserID:    c.UserID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}

func toDomain(m *models.TaskComment) *task_comments.TaskComment {
	return &task_comments.TaskComment{
		ID:            m.ID,
		TaskID:        m.TaskID,
		UserID:        m.UserID,
		UserFirstName: m.User.FirstName,
		UserLastName:  m.User.LastName,
		Content:       m.Content,
		CreatedAt:     m.CreatedAt,
	}
}
