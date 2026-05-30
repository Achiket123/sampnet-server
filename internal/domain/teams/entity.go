package teams

import (
	authDomain "server/internal/domain/auth"
	orgDomain "server/internal/domain/organisation"
)

type Team struct {
	ID             uint            `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	OrganisationID uint            `json:"organisation_id"`
	Organisation   orgDomain.Entity `json:"organisation"`
	CreatedBy      uint            `json:"created_by"`
	TeamLead       uint            `json:"team_lead"`
	TeamLeadUser   authDomain.User `json:"team_lead_user"`
	CreatedByUser  authDomain.User `json:"created_by_user"`
	IsActive       bool            `json:"is_active"`
}

type TeamMember struct {
	ID         uint            `json:"id"`
	UserID     uint            `json:"user_id"`
	TeamID     uint            `json:"team_id"`
	User       authDomain.User `json:"user"`
	Team       Team            `json:"team"`
	Role       string          `json:"role"`
	IsActive   bool            `json:"is_active"`
	IsLeader   bool            `json:"is_leader"`
	IsAdmin    bool            `json:"is_admin"`
	IsManager  bool            `json:"is_manager"`
	IsEmployee bool            `json:"is_employee"`
	IsBoss     bool            `json:"is_boss"`
}
