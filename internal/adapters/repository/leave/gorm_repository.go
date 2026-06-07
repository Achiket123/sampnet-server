package leave

import (
	"context"
	"server/internal/domain/leave"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) leave.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, l *leave.Leave) error {
	model := toModel(l)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	l.ID = model.ID
	l.CreatedAt = model.CreatedAt
	l.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*leave.Leave, error) {
	var model models.Leave
	if err := r.db.WithContext(ctx).Preload("Employee").Preload("Manager").First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) GetByEmployee(ctx context.Context, employeeID uint, offset int) ([]leave.Leave, error) {
	var modelsList []models.Leave
	if err := r.db.WithContext(ctx).
		Where("employee_id = ?", employeeID).
		Order("created_at DESC").
		Limit(20).
		Offset(offset * 20).
		Preload("Employee").
		Find(&modelsList).Error; err != nil {
		return nil, err
	}

	leaves := make([]leave.Leave, len(modelsList))
	for i, m := range modelsList {
		leaves[i] = *toDomain(&m)
	}
	return leaves, nil
}

func (r *gormRepository) GetByOrganisation(ctx context.Context, orgID uint, status string, offset int) ([]leave.Leave, error) {
	query := r.db.WithContext(ctx).Where("organisation_id = ?", orgID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var modelsList []models.Leave
	if err := query.
		Order("created_at DESC").
		Limit(20).
		Offset(offset * 20).
		Preload("Employee").
		Find(&modelsList).Error; err != nil {
		return nil, err
	}

	leaves := make([]leave.Leave, len(modelsList))
	for i, m := range modelsList {
		leaves[i] = *toDomain(&m)
	}
	return leaves, nil
}

func (r *gormRepository) GetHistoryByEmployee(ctx context.Context, employeeID uint, status string, from *time.Time, to *time.Time, limit int, offset int) ([]leave.Leave, error) {
	query := r.db.WithContext(ctx).Where("employee_id = ?", employeeID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if from != nil {
		query = query.Where("start_date >= ?", *from)
	}
	if to != nil {
		query = query.Where("end_date <= ?", *to)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	var modelsList []models.Leave
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Find(&modelsList).Error; err != nil {
		return nil, err
	}

	leaves := make([]leave.Leave, len(modelsList))
	for i, m := range modelsList {
		leaves[i] = *toDomain(&m)
	}
	return leaves, nil
}

func (r *gormRepository) Update(ctx context.Context, l *leave.Leave) error {
	model := toModel(l)
	return r.db.WithContext(ctx).Model(&models.Leave{}).Where("id = ?", l.ID).Save(model).Error
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Leave{}, id).Error
}

func toModel(l *leave.Leave) *models.Leave {
	var managerID *uint
	if l.ManagerID != 0 {
		managerID = &l.ManagerID
	}
	return &models.Leave{
		Model: gorm.Model{
			ID:        l.ID,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		},
		EmployeeID:     l.EmployeeID,
		OrganisationID: l.OrganisationID,
		ManagerID:      managerID,
		LeaveType:      l.LeaveType,
		StartDate:      l.StartDate,
		EndDate:        l.EndDate,
		TotalDays:      l.TotalDays,
		Reason:         l.Reason,
		Status:         l.Status,
		ManagerNote:    l.ManagerNote,
	}
}

func toDomain(m *models.Leave) *leave.Leave {
	var managerID uint
	if m.ManagerID != nil {
		managerID = *m.ManagerID
	}
	return &leave.Leave{
		ID:             m.ID,
		EmployeeID:     m.EmployeeID,
		OrganisationID: m.OrganisationID,
		ManagerID:      managerID,
		LeaveType:      m.LeaveType,
		StartDate:      m.StartDate,
		EndDate:        m.EndDate,
		TotalDays:      m.TotalDays,
		Reason:         m.Reason,
		Status:         m.Status,
		ManagerNote:    m.ManagerNote,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
