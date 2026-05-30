package projects

import (
	"time"
	authDomain "server/internal/domain/auth"
	orgDomain "server/internal/domain/organisation"
	teamDomain "server/internal/domain/teams"
)

type Project struct {
	ID               uint               `json:"id"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	StartDate        time.Time          `json:"start_date"`
	EndDate          time.Time          `json:"end_date"`
	OrganisationID   uint               `json:"organisation_id"`
	Organisation     orgDomain.Entity    `json:"organisation"`
	TeamID           uint               `json:"team_id"`
	Team             teamDomain.Team    `json:"team"`
	CreatedBy        uint               `json:"created_by"`
	CreatedByUser    authDomain.User    `json:"created_by_user"`
	Status           string             `json:"status"`
	Priority         string             `json:"priority"`
	CompletionStatus string             `json:"completion_status"`
}
