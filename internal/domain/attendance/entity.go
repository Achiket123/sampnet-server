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
}
