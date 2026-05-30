package employees

import (
	"time"
	authDomain "server/internal/domain/auth"
	orgDomain "server/internal/domain/organisation"
)

type Employee struct {
	UserID         uint                   `json:"user_id"`
	User           authDomain.User        `json:"user"`
	EmploymentID   int                    `json:"employment_id"`
	OrganisationID uint                   `json:"organisation_id"`
	Organisation   orgDomain.Entity       `json:"organisation"`
	Type           string                 `json:"type"`
	Salary         string                 `json:"salary"`
	LastLoginAt    time.Time              `json:"last_login_at"`
}

type Manager struct {
	UserID         uint                   `json:"user_id"`
	User           authDomain.User        `json:"user"`
	EmploymentID   int                    `json:"employment_id"`
	OrganisationID uint                   `json:"organisation_id"`
	Organisation   orgDomain.Entity       `json:"organisation"`
	Type           string                 `json:"type"`
	Salary         string                 `json:"salary"`
	LastLoginAt    time.Time              `json:"last_login_at"`
}

type Boss struct {
	UserID         uint                   `json:"user_id"`
	User           authDomain.User        `json:"user"`
	OrganisationID uint                   `json:"organisation_id"`
	Organisation   orgDomain.Entity       `json:"organisation"`
	LastLoginAt    time.Time              `json:"last_login_at"`
}
