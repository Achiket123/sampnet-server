package leave

import (
	"errors"
	"time"
)

var (
	ErrLeaveNotFound = errors.New("leave request not found")
	ErrNotAuthorized = errors.New("not authorized to perform this action")
)

type Leave struct {
	ID             uint      `json:"id"`
	EmployeeID     uint      `json:"employee_id"`
	OrganisationID uint      `json:"organisation_id"`
	ManagerID      uint      `json:"manager_id"`
	LeaveType      string    `json:"leave_type"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	TotalDays      int       `json:"total_days"`
	Reason         string    `json:"reason"`
	Status         string    `json:"status"`
	ManagerNote    string    `json:"manager_note"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LeaveSummary struct {
	TotalLeavesTaken int            `json:"total_leaves_taken"`
	ApprovedCount    int            `json:"approved_count"`
	PendingCount     int            `json:"pending_count"`
	RejectedCount    int            `json:"rejected_count"`
	LeavesByType     map[string]int `json:"leaves_by_type"`
}

type LeaveHistory struct {
	Records []Leave      `json:"records"`
	Summary LeaveSummary `json:"summary"`
}
