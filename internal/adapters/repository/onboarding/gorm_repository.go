package onboarding

import (
	"context"
	"server/internal/domain/onboarding"
	"server/internal/platform/database/models"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) onboarding.Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetByUserID(ctx context.Context, userID uint) (*onboarding.OnboardingProgress, error) {
	var m models.OnboardingProgress
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &onboarding.OnboardingProgress{UserID: userID}, nil
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *gormRepository) Update(ctx context.Context, p *onboarding.OnboardingProgress) error {
	var m models.OnboardingProgress
	err := r.db.WithContext(ctx).Where("user_id = ?", p.UserID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			m = *toModel(p)
			err = r.db.WithContext(ctx).Create(&m).Error
			if err == nil {
				if m.ProfileCompleted && m.OrganisationID != 0 && m.TeamJoined && m.TaskCreated && m.InviteSent {
					r.db.WithContext(ctx).Model(&models.Employee{}).Where("user_id = ?", p.UserID).Update("onboarding_completed", true)
				}
			}
			return err
		}
		return err
	}

	m.OrganisationID = p.OrganisationID
	m.ProfileCompleted = p.ProfileCompleted
	m.TeamJoined = p.TeamJoined
	m.TaskCreated = p.TaskCreated
	m.InviteSent = p.InviteSent
	
	err = r.db.WithContext(ctx).Save(&m).Error
	if err == nil {
		if m.ProfileCompleted && m.OrganisationID != 0 && m.TeamJoined && m.TaskCreated && m.InviteSent {
			r.db.WithContext(ctx).Model(&models.Employee{}).Where("user_id = ?", p.UserID).Update("onboarding_completed", true)
		}
	}
	return err
}

func toModel(p *onboarding.OnboardingProgress) *models.OnboardingProgress {
	return &models.OnboardingProgress{
		UserID:           p.UserID,
		OrganisationID:   p.OrganisationID,
		ProfileCompleted: p.ProfileCompleted,
		TeamJoined:       p.TeamJoined,
		TaskCreated:      p.TaskCreated,
		InviteSent:       p.InviteSent,
	}
}

func toDomain(m *models.OnboardingProgress) *onboarding.OnboardingProgress {
	return &onboarding.OnboardingProgress{
		UserID:           m.UserID,
		OrganisationID:   m.OrganisationID,
		ProfileCompleted: m.ProfileCompleted,
		TeamJoined:       m.TeamJoined,
		TaskCreated:      m.TaskCreated,
		InviteSent:       m.InviteSent,
	}
}
