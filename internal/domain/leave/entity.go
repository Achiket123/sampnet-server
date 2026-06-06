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
