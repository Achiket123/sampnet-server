package analytics

import (
	"context"
	"server/internal/domain/employees"
	"time"
)

type AnalyticsAttendanceBlock struct {
	PresentDays            int     `json:"present_days"`
	AbsentDays             int     `json:"absent_days"`
	TotalWorkingDays       int     `json:"total_working_days"`
	AttendanceRate         float64 `json:"attendance_rate"`
	AverageCheckInTime     string  `json:"average_check_in_time"`
	AverageDurationMinutes int     `json:"average_duration_minutes"`
	LateCheckins           int     `json:"late_checkins"`
}

type AnalyticsLeaveBlock struct {
	TotalLeavesInPeriod int            `json:"total_leaves_in_period"`
	Approved            int            `json:"approved"`
	Pending             int            `json:"pending"`
	Rejected            int            `json:"rejected"`
	LeavesByType        map[string]int `json:"leaves_by_type"`
}

type AnalyticsTaskBlock struct {
	TotalAssigned  int     `json:"total_assigned"`
	Completed      int     `json:"completed"`
	InProgress     int     `json:"in_progress"`
	Overdue        int     `json:"overdue"`
	CompletionRate float64 `json:"completion_rate"`
}

type EmployeeAnalyticsSummary struct {
	EmployeeInfo employees.Employee       `json:"employee_info"`
	Attendance   AnalyticsAttendanceBlock `json:"attendance"`
	Leave        AnalyticsLeaveBlock      `json:"leave"`
	Tasks        AnalyticsTaskBlock       `json:"tasks"`
}

type EmployeeMonitorSnapshot struct {
	AttendanceRateThisMonth float64    `json:"attendance_rate_this_month"`
	LeavesPending           int        `json:"leaves_pending"`
	TasksOverdue            int        `json:"tasks_overdue"`
	LastSeen                *time.Time `json:"last_seen"`
	Status                  string     `json:"status"` // active, idle, inactive
}

type EmployeeMonitorResponse struct {
	EmployeeInfo employees.Employee      `json:"employee_info"`
	Snapshot     EmployeeMonitorSnapshot `json:"snapshot"`
}

type Repository interface {
	GetEmployeeAnalytics(ctx context.Context, userID uint, orgID uint, period string) (*EmployeeAnalyticsSummary, error)
	GetOrgEmployeeMonitor(ctx context.Context, orgID uint) ([]EmployeeMonitorResponse, error)
}

type UseCase interface {
	GetEmployeeAnalytics(ctx context.Context, userID uint, orgID uint, period string) (*EmployeeAnalyticsSummary, error)
	GetOrgEmployeeMonitor(ctx context.Context, orgID uint) ([]EmployeeMonitorResponse, error)
}