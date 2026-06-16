package invites

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	authDomain "server/internal/domain/auth"
	employeesDomain "server/internal/domain/employees"
	domain "server/internal/domain/invites"
	"server/internal/platform/database/models"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) domain.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, invite *domain.EmployeeInvite) error {
	m := toModel(invite)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	invite.ID = m.ID
	invite.CreatedAt = m.CreatedAt
	return nil
}

func (r *gormRepository) GetByToken(ctx context.Context, token string) (*domain.EmployeeInvite, error) {
	var m models.EmployeeInvite
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invite not found")
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *gormRepository) GetByEmail(ctx context.Context, email string, orgID uint) (*domain.EmployeeInvite, error) {
	var m models.EmployeeInvite
	if err := r.db.WithContext(ctx).Where("email = ? AND organisation_id = ? AND status = ?", email, orgID, "pending").First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invite not found")
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *gormRepository) MarkAccepted(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.EmployeeInvite{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "accepted",
			"accepted_at": time.Now(),
		}).Error
}

func (r *gormRepository) GetByOrg(ctx context.Context, orgID uint) ([]domain.EmployeeInvite, error) {
	var ms []models.EmployeeInvite
	if err := r.db.WithContext(ctx).Where("organisation_id = ?", orgID).Order("created_at DESC").Find(&ms).Error; err != nil {
		return nil, err
	}
	res := make([]domain.EmployeeInvite, len(ms))
	for i, m := range ms {
		res[i] = *toDomain(&m)
	}
	return res, nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*domain.EmployeeInvite, error) {
	var m models.EmployeeInvite
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invite not found")
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *gormRepository) Update(ctx context.Context, invite *domain.EmployeeInvite) error {
	m := toModel(invite)
	return r.db.WithContext(ctx).Save(m).Error
}


func (r *gormRepository) AcceptInvite(ctx context.Context, inviteID uint, u *authDomain.User, emp *employeesDomain.Employee) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userModel := &models.UserModel{
			FirstName:      u.FirstName,
			LastName:       u.LastName,
			Email:          u.Email,
			PhoneNumber:    u.PhoneNumber,
			IsVerified:     u.IsVerified,
			HashedPassword: u.HashedPassword,
		}
		userModel.ID = u.ID
		if err := tx.Save(userModel).Error; err != nil {
			return err
		}
		u.ID = userModel.ID
		emp.UserID = userModel.ID

		employeeModel := &models.Employee{
			UserID:         emp.UserID,
			EmploymentID:   emp.EmploymentID,
			OrganisationID: emp.OrganisationID,
			Email:          emp.Email,
		}
		if err := tx.Create(employeeModel).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.EmployeeInvite{}).
			Where("id = ?", inviteID).
			Updates(map[string]interface{}{
				"status":      "accepted",
				"accepted_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func toModel(invite *domain.EmployeeInvite) *models.EmployeeInvite {
	m := &models.EmployeeInvite{
		Token:           invite.Token,
		Email:           invite.Email,
		FirstName:       invite.FirstName,
		LastName:        invite.LastName,
		PhoneNumber:     invite.PhoneNumber,
		EmploymentID:    invite.EmploymentID,
		OrganisationID:  invite.OrganisationID,
		InvitedByUserID: invite.InvitedByUserID,
		Status:          invite.Status,
		ExpiresAt:       invite.ExpiresAt,
		AcceptedAt:      invite.AcceptedAt,
	}
	m.ID = invite.ID
	return m
}

func toDomain(m *models.EmployeeInvite) *domain.EmployeeInvite {
	return &domain.EmployeeInvite{
		ID:              m.ID,
		Token:           m.Token,
		Email:           m.Email,
		FirstName:       m.FirstName,
		LastName:        m.LastName,
		PhoneNumber:     m.PhoneNumber,
		EmploymentID:    m.EmploymentID,
		OrganisationID:  m.OrganisationID,
		InvitedByUserID: m.InvitedByUserID,
		Status:          m.Status,
		ExpiresAt:       m.ExpiresAt,
		CreatedAt:       m.CreatedAt,
		AcceptedAt:      m.AcceptedAt,
	}
}