package organisation

import (
	"context"
	domain "server/internal/domain/organisation"
	"server/internal/platform/database/models"
	"time"

	"gorm.io/gorm"
)

// GormRepository implements organisation.Repository with GORM.
type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, org *domain.Entity) error {
	model := toModel(org)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	org.ID = model.ID
	return nil
}

func (r *GormRepository) CreateWithOwner(ctx context.Context, org *domain.Entity, ownerUserID uint) (*domain.OwnerEmployeeRow, error) {
	var createdEmployee domain.OwnerEmployeeRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := toModel(org)
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		org.ID = model.ID

		employee := models.Employee{
			UserID:         ownerUserID,
			OrganisationID: model.ID,
			EmploymentID:   int(model.ID),
			Type:           "owner",
			LastLoginAt:    time.Now().UTC(),
		}

		if err := tx.Create(&employee).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Organisation{}).Where("id = ?", model.ID).Update("boss_id", ownerUserID).Error; err != nil {
			return err
		}

		org.BossID = ownerUserID
		createdEmployee = domain.OwnerEmployeeRow{
			UserID:         employee.UserID,
			EmploymentID:   employee.EmploymentID,
			OrganisationID: employee.OrganisationID,
			Type:           employee.Type,
			Salary:         employee.Salary,
			LastLoginAt:    employee.LastLoginAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &createdEmployee, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (*domain.Entity, error) {
	var model models.Organisation
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	entity := toEntity(model)
	return &entity, nil
}

func (r *GormRepository) Update(ctx context.Context, org *domain.Entity) error {
	model := toModel(org)
	return r.db.WithContext(ctx).Model(&models.Organisation{}).Where("id = ?", org.ID).Updates(&model).Error
}

func toModel(org *domain.Entity) models.Organisation {
	return models.Organisation{
		Model:              gorm.Model{ID: org.ID},
		CompanyName:        org.CompanyName,
		CompanyCode:        org.CompanyCode,
		PrimaryContactName: org.PrimaryContactName,
		PrimaryEmail:       org.PrimaryEmail,
		PhoneNumber:        org.PhoneNumber,
		OfficeAddress:      org.OfficeAddress,
		City:               org.City,
		State:              org.State,
		PostalCode:         org.PostalCode,
		Country:            org.Country,
		PlanID:             org.PlanID,
		PlanStartDate:      org.PlanStartDate,
		PlanEndDate:        org.PlanEndDate,
		PlanStatus:         org.PlanStatus,
		MaxEmployees:       org.MaxEmployees,
		CompanyLogo:        org.CompanyLogo,
		Industry:           org.Industry,
		BillingAddress:     org.BillingAddress,
		CompanySize:        org.CompanySize,
		BossID:             org.BossID,
	}
}

func toEntity(model models.Organisation) domain.Entity {
	return domain.Entity{
		ID:                 model.ID,
		CompanyName:        model.CompanyName,
		CompanyCode:        model.CompanyCode,
		PrimaryContactName: model.PrimaryContactName,
		PrimaryEmail:       model.PrimaryEmail,
		PhoneNumber:        model.PhoneNumber,
		OfficeAddress:      model.OfficeAddress,
		City:               model.City,
		State:              model.State,
		PostalCode:         model.PostalCode,
		Country:            model.Country,
		PlanID:             model.PlanID,
		PlanStartDate:      model.PlanStartDate,
		PlanEndDate:        model.PlanEndDate,
		PlanStatus:         model.PlanStatus,
		MaxEmployees:       model.MaxEmployees,
		CompanyLogo:        model.CompanyLogo,
		Industry:           model.Industry,
		BillingAddress:     model.BillingAddress,
		CompanySize:        model.CompanySize,
		BossID:             model.BossID,
	}
}
