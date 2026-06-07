package attendance

import (
	"time"
)

type Attendance struct {
	ID             uint       `json:"id"`
	UserID         uint       `json:"user_id"`
	Date           time.Time  `json:"date"`
	CheckInTime    *time.Time `json:"check_in_time"`
	CheckOutTime   *time.Time `json:"check_out_time"`
	OrganisationID uint       `json:"organisation_id"`
	CheckInPhoto   int        `json:"check_in_photo"`
	CheckOutPhoto  int        `json:"check_out_photo"`
	DurationMinutes *int      `json:"duration_minutes"`
}

type AttendanceSummary struct {
	TotalDaysPresent       int     `json:"total_days_present"`
	TotalDaysAbsent        int     `json:"total_days_absent"`
	AverageCheckInTime     string  `json:"average_check_in_time"`
	AverageDurationMinutes int     `json:"average_duration_minutes"`
	OnTimePercentage       float64 `json:"on_time_percentage"`
}

type AttendanceHistory struct {
	Records []Attendance      `json:"records"`
	Summary AttendanceSummary `json:"summary"`
}
