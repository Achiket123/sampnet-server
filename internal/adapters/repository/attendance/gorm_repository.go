package attendance

import (
	"context"
	domain "server/internal/domain/attendance"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, att *domain.Attendance) error {
	model := toModel(att)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	att.ID = model.ID
	return nil
}

func (r *gormRepository) GetByUserAndDate(ctx context.Context, userID uint, date time.Time) (*domain.Attendance, error) {
	var model models.Attendance
	if err := r.db.WithContext(ctx).Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).First(&model).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) Update(ctx context.Context, att *domain.Attendance) error {
	model := toModel(att)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *gormRepository) GetByUser(ctx context.Context, userID uint, offset int) ([]domain.Attendance, error) {
	var modelsList []models.Attendance
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func (r *gormRepository) GetByOrganisation(ctx context.Context, orgID uint, offset int) ([]domain.Attendance, error) {
	var modelsList []models.Attendance
	if err := r.db.WithContext(ctx).Where("organisation_id = ?", orgID).Offset(offset).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	return toDomainList(modelsList), nil
}

func toModel(a *domain.Attendance) *models.Attendance {
	model := &models.Attendance{
		UserID:         a.UserID,
		Date:           a.Date,
		CheckInTime:    a.CheckInTime,
		CheckOutTime:   a.CheckOutTime,
		OrganisationID: a.OrganisationID,
		CheckInPhoto:   a.CheckInPhoto,
		CheckOutPhoto:  a.CheckOutPhoto,
	}
	model.ID = a.ID
	return model
}

func toDomain(m *models.Attendance) *domain.Attendance {
	return &domain.Attendance{
		ID:             m.ID,
		UserID:         m.UserID,
		Date:           m.Date,
		CheckInTime:    m.CheckInTime,
		CheckOutTime:   m.CheckOutTime,
		OrganisationID: m.OrganisationID,
		CheckInPhoto:   m.CheckInPhoto,
		CheckOutPhoto:  m.CheckOutPhoto,
	}
}

func toDomainList(modelsList []models.Attendance) []domain.Attendance {
	result := make([]domain.Attendance, len(modelsList))
	for i, m := range modelsList {
		result[i] = *toDomain(&m)
	}
	return result
}
