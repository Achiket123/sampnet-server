package projects

import (
	"encoding/json"
	"fmt"
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
	Milestones       []Milestone        `json:"milestones"`
}

type Milestone struct {
	ID             uint      `json:"id"`
	ProjectID      uint      `json:"project_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	DueDate        time.Time `json:"due_date"`
	Status         string    `json:"status"` // e.g. "Pending", "Completed"
	IsOverdue      bool      `json:"is_overdue"`
	OrganisationID uint      `json:"organisation_id"`
}

func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	// Try RFC3339 first
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// Try without Z (ISO 8601 local time format with milliseconds)
	if t, err := time.Parse("2006-01-02T15:04:05.999999999", s); err == nil {
		return t, nil
	}
	// Try ISO 8601 local time format without milliseconds
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		return t, nil
	}
	// Try date only
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", s)
}

func (p *Project) UnmarshalJSON(data []byte) error {
	type Alias Project
	aux := &struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.StartDate != "" {
		t, err := parseTime(aux.StartDate)
		if err != nil {
			return err
		}
		p.StartDate = t
	}
	if aux.EndDate != "" {
		t, err := parseTime(aux.EndDate)
		if err != nil {
			return err
		}
		p.EndDate = t
	}
	return nil
}

func (m *Milestone) UnmarshalJSON(data []byte) error {
	type Alias Milestone
	aux := &struct {
		DueDate string `json:"due_date"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.DueDate != "" {
		t, err := parseTime(aux.DueDate)
		if err != nil {
			return err
		}
		m.DueDate = t
	}
	return nil
}
