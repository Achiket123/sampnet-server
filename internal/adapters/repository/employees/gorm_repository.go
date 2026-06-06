package employees

import (
	"context"
	authDomain "server/internal/domain/auth"
	domain "server/internal/domain/employees"
	orgDomain "server/internal/domain/organisation"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetByUserID(ctx context.Context, userID uint) (*domain.Employee, error) {
	var model models.Employee
	if err := r.db.WithContext(ctx).Preload("User").Preload("Organisation").Where("user_id = ?", userID).First(&model).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) Create(ctx context.Context, emp *domain.Employee) error {
	model := toModel(emp)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *gormRepository) GetEmployeesByOrg(ctx context.Context, orgID uint) ([]domain.Employee, error) {
	var modelsList []models.Employee
	if err := r.db.WithContext(ctx).Preload("User").Preload("Organisation").Where("organisation_id = ?", orgID).Find(&modelsList).Error; err != nil {
		return nil, err
	}
	result := make([]domain.Employee, len(modelsList))
	for i, m := range modelsList {
		result[i] = *toDomain(&m)
	}
	return result, nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*domain.Employee, error) {
	var model models.Employee
	if err := r.db.WithContext(ctx).Preload("User").Preload("Organisation").First(&model, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) Update(ctx context.Context, emp *domain.Employee) error {
	model := toModel(emp)
	return r.db.WithContext(ctx).Model(&models.Employee{}).Where("user_id = ?", emp.UserID).Updates(model).Error
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", id).Delete(&models.Employee{}).Error
}

func (r *gormRepository) Search(ctx context.Context, query string) ([]domain.Employee, error) {
	var modelsList []models.Employee
	// Searching in UserModel through joining
	if err := r.db.WithContext(ctx).Preload("User").Preload("Organisation").
		Joins("JOIN user_models ON user_models.id = employees.user_id").
		Where("user_models.first_name LIKE ? OR user_models.last_name LIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&modelsList).Error; err != nil {
		return nil, err
	}
	result := make([]domain.Employee, len(modelsList))
	for i, m := range modelsList {
		result[i] = *toDomain(&m)
	}
	return result, nil
}

func (r *gormRepository) CreateManager(ctx context.Context, manager *domain.Manager) error {
	model := &models.Manager{
		UserID:         manager.UserID,
		OrganisationID: manager.OrganisationID,
		Type:           manager.Type,
		Email:          manager.Email,
		Salary:         manager.Salary,
		LastLoginAt:    manager.LastLoginAt,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *gormRepository) GetManagerByUserID(ctx context.Context, userID uint) (*domain.Manager, error) {
	var model models.Manager
	if err := r.db.WithContext(ctx).Preload("User").Preload("Organisation").Where("user_id = ?", userID).First(&model).Error; err != nil {
		return nil, err
	}
	return &domain.Manager{
		UserID: model.UserID,
		User: authDomain.User{
			ID:          model.User.ID,
			FirstName:   model.User.FirstName,
			LastName:    model.User.LastName,
			Email:       model.User.Email,
			PhoneNumber: model.User.PhoneNumber,
		},
		OrganisationID: model.OrganisationID,
		Organisation: orgDomain.Entity{
			ID:          model.Organisation.ID,
			CompanyName: model.Organisation.CompanyName,
		},
		Email:       model.Email,
		Type:        model.Type,
		Salary:      model.Salary,
		LastLoginAt: model.LastLoginAt,
	}, nil
}

func (r *gormRepository) CreateBoss(ctx context.Context, boss *domain.Boss) error {
	model := &models.Boss{
		UserID:         boss.UserID,
		OrganisationID: boss.OrganisationID,
		LastLoginAt:    boss.LastLoginAt,
		Email:          boss.Email,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *gormRepository) GetBossByUserID(ctx context.Context, userID uint) (*domain.Boss, error) {
	var model models.Boss
	if err := r.db.WithContext(ctx).Preload("User").Preload("Organisation").Where("user_id = ?", userID).First(&model).Error; err != nil {
		return nil, err
	}
	return &domain.Boss{
		UserID: model.UserID,
		User: authDomain.User{
			ID:          model.User.ID,
			FirstName:   model.User.FirstName,
			LastName:    model.User.LastName,
			Email:       model.User.Email,
			PhoneNumber: model.User.PhoneNumber,
		},
		Email:          model.Email,
		OrganisationID: model.OrganisationID,
		Organisation: orgDomain.Entity{
			ID:          model.Organisation.ID,
			CompanyName: model.Organisation.CompanyName,
		},
		LastLoginAt: model.LastLoginAt,
	}, nil
}

func toModel(e *domain.Employee) *models.Employee {
	return &models.Employee{
		UserID:         e.UserID,
		EmploymentID:   e.EmploymentID,
		OrganisationID: e.OrganisationID,
		Email:          e.Email,
		Type:           e.Type,
		Salary:         e.Salary,
		LastLoginAt:    e.LastLoginAt,
	}
}

func toDomain(m *models.Employee) *domain.Employee {
	return &domain.Employee{
		UserID: m.UserID,
		User: authDomain.User{
			ID:          m.User.ID,
			FirstName:   m.User.FirstName,
			LastName:    m.User.LastName,
			Email:       m.User.Email,
			PhoneNumber: m.User.PhoneNumber,
		},
		EmploymentID:   m.EmploymentID,
		OrganisationID: m.OrganisationID,
		Organisation: orgDomain.Entity{
			ID:          m.Organisation.ID,
			CompanyName: m.Organisation.CompanyName,
		},
		Email:       m.Email,
		Type:        m.Type,
		Salary:      m.Salary,
		LastLoginAt: m.LastLoginAt,
	}
}
