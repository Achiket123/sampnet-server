package tasks

import (
	authDomain "server/internal/domain/auth"
	orgDomain "server/internal/domain/organisation"
	projectDomain "server/internal/domain/projects"
	teamDomain "server/internal/domain/teams"
	"time"
)

type Task struct {
	ID             uint                   `json:"id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	DueDate        time.Time              `json:"due_date"`
	AssignedTo     uint                   `json:"assigned_to"`
	AssignedUser   *authDomain.User       `json:"assigned_user"`
	AssignedBy     uint                   `json:"assigned_by"`
	AssignedByUser *authDomain.User       `json:"assigned_by_user"`
	Type           string                 `json:"type"`
	Priority       string                 `json:"priority"`
	Status         string                 `json:"status"`
	OrganisationID uint                   `json:"organisation_id"`
	Organisation   *orgDomain.Entity      `json:"organisation"`
	IsPersonal     bool                   `json:"is_personal"`
	TeamID         *uint                  `json:"team_id"`
	Team           *teamDomain.Team       `json:"team"`
	ProjectID      *uint                  `json:"project_id"`
	Project        *projectDomain.Project `json:"project"`
	CreatedAt      time.Time              `json:"CreatedAt"`
}
