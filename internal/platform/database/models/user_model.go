package models

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	FirstName      string    `gorm:"column:first_name"  json:"first_name" validate:"required,min=2,max=100"`
	LastName       string    `gorm:"column:last_name" json:"last_name" validate:"required,min=2,max=100"`
	Email          string    `gorm:"column:email;unique;Not Null" json:"email" validate:"required,email"`
	PhoneNumber    string    `gorm:"column:phone_number;unique" json:"phone_number" validate:"required,len=10"`
	IsVerified     bool      `gorm:"column:is_verified;default:false" json:"is_verified"`
	HashedPassword string    `gorm:"column:hashed_password" `
	ProfilePic     string    `gorm:"column:profile_id" json:"profile_id"`
	City           string    `gorm:"column:city" json:"city"`
	Country        string    `gorm:"column:country" json:"country"`
	DateOfBirth    time.Time `gorm:"column:date_of_birth" json:"date_of_birth"`
	LastLoginAt    time.Time `gorm:"column:last_login_at" json:"last_login_at"`
}

type Employee struct {
	UserID         uint           `gorm:"column:user_id;primaryKey" json:"user_id"`
	User           UserModel      `gorm:"foreignKey:UserID;references:ID"`
	EmploymentID   int            `gorm:"column:employment_id" json:"employment_id"`
	OrganisationID uint           `gorm:"column:organisation_id;references:ID" json:"organisation_id"`
	Organisation   Organisation   `gorm:"foreignKey:OrganisationID;references:ID"`
	Type           string         `gorm:"column:type" json:"type"`
	Salary         string         `gorm:"column:salary" json:"salary"`
	Email          string         `gorm:"column:email" json:"email"`
	LastLoginAt    time.Time      `gorm:"column:last_login_at" json:"last_login_at"`
	OnboardingCompleted bool      `gorm:"column:onboarding_completed;default:false" json:"onboarding_completed"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Manager struct {
	UserID         uint           `gorm:"column:user_id;primaryKey" json:"user_id"`
	User           UserModel      `gorm:"foreignKey:UserID;references:ID"`
	EmploymentID   int            `gorm:"column:employment_id" json:"employment_id"`
	OrganisationID uint           `gorm:"column:organisation_id" json:"organisation_id"`
	Organisation   Organisation   `gorm:"foreignKey:OrganisationID;references:ID"`
	Type           string         `gorm:"column:type" json:"type"`
	Salary         string         `gorm:"column:salary" json:"salary"`
	Email          string         `gorm:"column:email" json:"email"`
	LastLoginAt    time.Time      `gorm:"column:last_login_at" json:"last_login_at"`
	OnboardingCompleted bool      `gorm:"column:onboarding_completed;default:false" json:"onboarding_completed"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Boss struct {
	UserID         uint           `gorm:"column:user_id;primaryKey" json:"user_id"`
	User           UserModel      `gorm:"foreignKey:UserID;references:ID"`
	OrganisationID uint           `gorm:"column:organisation_id" json:"organisation_id"`
	Organisation   Organisation   `gorm:"foreignKey:OrganisationID;references:ID"`
	Email          string         `gorm:"column:email" json:"email"`
	LastLoginAt    time.Time      `gorm:"column:last_login_at" json:"last_login_at"`
	OnboardingCompleted bool      `gorm:"column:onboarding_completed;default:false" json:"onboarding_completed"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type EmployeeInvite struct {
	gorm.Model
	Token           string     `gorm:"column:token;uniqueIndex;not null" json:"token"`
	Email           string     `gorm:"column:email;not null" json:"email"`
	FirstName       string     `gorm:"column:first_name" json:"first_name"`
	LastName        string     `gorm:"column:last_name" json:"last_name"`
	PhoneNumber     string     `gorm:"column:phone_number" json:"phone_number"`
	EmploymentID    int        `gorm:"column:employment_id" json:"employment_id"`
	OrganisationID  uint       `gorm:"column:organisation_id;not null" json:"organisation_id"`
	InvitedByUserID uint       `gorm:"column:invited_by_user_id;not null" json:"invited_by_user_id"`
	Status          string     `gorm:"column:status;default:'pending'" json:"status"`
	ExpiresAt       time.Time  `gorm:"column:expires_at" json:"expires_at"`
	AcceptedAt      *time.Time `gorm:"column:accepted_at" json:"accepted_at"`
}

type EmailVerification struct {
	gorm.Model
	UserID    uint       `gorm:"column:user_id;uniqueIndex;not null" json:"user_id"`
	Token     string     `gorm:"column:token;uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null" json:"expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"used_at"`
}


