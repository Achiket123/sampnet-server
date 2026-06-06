package employees

import (
	authDomain "server/internal/domain/auth"
	orgDomain "server/internal/domain/organisation"
	"time"
)

type Employee struct {
	UserID         uint             `json:"user_id"`
	User           authDomain.User  `json:"user"`
	EmploymentID   int              `json:"employment_id"`
	OrganisationID uint             `json:"organisation_id"`
	Email          string           `json:"email"`
	Organisation   orgDomain.Entity `json:"organisation"`
	Type           string           `json:"type"`
	Salary         string           `json:"salary"`
	LastLoginAt    time.Time        `json:"last_login_at"`
	OnboardingCompleted bool        `json:"onboarding_completed"`
}

type Manager struct {
	UserID         uint             `json:"user_id"`
	User           authDomain.User  `json:"user"`
	EmploymentID   int              `json:"employment_id"`
	OrganisationID uint             `json:"organisation_id"`
	Email          string           `json:"email"`
	Organisation   orgDomain.Entity `json:"organisation"`

	Type        string    `json:"type"`
	Salary      string    `json:"salary"`
	LastLoginAt time.Time `json:"last_login_at"`
	OnboardingCompleted bool        `json:"onboarding_completed"`
}

type Boss struct {
	UserID         uint             `json:"user_id"`
	User           authDomain.User  `json:"user"`
	Email          string           `json:"email"`
	OrganisationID uint             `json:"organisation_id"`
	Organisation   orgDomain.Entity `json:"organisation"`
	LastLoginAt    time.Time        `json:"last_login_at"`
	OnboardingCompleted bool        `json:"onboarding_completed"`
}
