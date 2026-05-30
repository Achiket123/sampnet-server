package attendance

import (
	"context"
	"time"
)

type UseCase interface {
	RecordAttendance(ctx context.Context, att *Attendance) error
	UpdateAttendance(ctx context.Context, userID uint, checkOutTime *time.Time, checkOutPhoto int) (*Attendance, error)
	GetAttendanceByUser(ctx context.Context, userID uint, offset int) ([]Attendance, error)
	GetAttendanceByOrganisation(ctx context.Context, orgID uint, offset int) ([]Attendance, error)
	GetAttendanceByDateAndUser(ctx context.Context, userID uint, date time.Time) (*Attendance, error)
}
