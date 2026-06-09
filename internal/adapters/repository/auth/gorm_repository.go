package auth

import (
	"context"
	domain "server/internal/domain/auth"
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

func (r *gormRepository) Create(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	user.ID = model.ID
	return nil
}

func (r *gormRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) GetByPhoneNumber(ctx context.Context, phone string) (*domain.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("phone_number = ?", phone).First(&model).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *gormRepository) Update(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *gormRepository) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	model := &models.RefreshToken{
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
		Revoked:   false,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *gormRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	var model models.RefreshToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&model).Error; err != nil {
		return nil, err
	}
	return &domain.RefreshToken{
		ID:        model.ID,
		UserID:    model.UserID,
		TokenHash: model.TokenHash,
		ExpiresAt: model.ExpiresAt,
		Revoked:   model.Revoked,
	}, nil
}

func (r *gormRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).Error
}

func (r *gormRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func toModel(u *domain.User) *models.UserModel {
	model := &models.UserModel{
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          u.Email,
		PhoneNumber:    u.PhoneNumber,
		IsVerified:     u.IsVerified,
		HashedPassword: u.HashedPassword,
		ProfilePic:     u.ProfilePic,
		City:           u.City,
		Country:        u.Country,
		DateOfBirth:    u.DateOfBirth,
		LastLoginAt:    u.LastLoginAt,
	}
	model.ID = u.ID
	return model
}

func toDomain(m *models.UserModel) *domain.User {
	return &domain.User{
		ID:             m.ID,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		Email:          m.Email,
		PhoneNumber:    m.PhoneNumber,
		IsVerified:     m.IsVerified,
		HashedPassword: m.HashedPassword,
		ProfilePic:     m.ProfilePic,
		City:           m.City,
		Country:        m.Country,
		DateOfBirth:    m.DateOfBirth,
		LastLoginAt:    m.LastLoginAt,
	}
}

func (r *gormRepository) CreateEmailVerification(ctx context.Context, ev *domain.EmailVerification) error {
	model := toEvModel(ev)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	ev.ID = model.ID
	return nil
}

func (r *gormRepository) GetEmailVerificationByToken(ctx context.Context, token string) (*domain.EmailVerification, error) {
	var model models.EmailVerification
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&model).Error; err != nil {
		return nil, err
	}
	return toEvDomain(&model), nil
}

func (r *gormRepository) GetActiveEmailVerificationByUserID(ctx context.Context, userID uint) (*domain.EmailVerification, error) {
	var model models.EmailVerification
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ? AND used_at IS NULL", userID, time.Now()).
		First(&model).Error; err != nil {
		return nil, err
	}
	return toEvDomain(&model), nil
}

func (r *gormRepository) UpdateEmailVerification(ctx context.Context, ev *domain.EmailVerification) error {
	model := toEvModel(ev)
	return r.db.WithContext(ctx).Save(model).Error
}

func toEvModel(ev *domain.EmailVerification) *models.EmailVerification {
	model := &models.EmailVerification{
		UserID:    ev.UserID,
		Token:     ev.Token,
		ExpiresAt: ev.ExpiresAt,
		UsedAt:    ev.UsedAt,
	}
	model.ID = ev.ID
	return model
}

func toEvDomain(m *models.EmailVerification) *domain.EmailVerification {
	return &domain.EmailVerification{
		ID:        m.ID,
		UserID:    m.UserID,
		Token:     m.Token,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
	}
}

