package auth

import (
	"context"
	domain "server/internal/domain/auth"
	"server/internal/platform/database/models"

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
