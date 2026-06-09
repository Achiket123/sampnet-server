package invites

import "time"

type EmployeeInvite struct {
	ID              uint       `json:"id"`
	Token           string     `json:"token"`
	Email           string     `json:"email"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	PhoneNumber     string     `json:"phone_number"`
	EmploymentID    int        `json:"employment_id"`
	OrganisationID  uint       `json:"organisation_id"`
	InvitedByUserID uint       `json:"invited_by_user_id"`
	Status          string     `json:"status"`
	ExpiresAt       time.Time  `json:"expires_at"`
	CreatedAt       time.Time  `json:"created_at"`
	AcceptedAt      *time.Time `json:"accepted_at"`
}