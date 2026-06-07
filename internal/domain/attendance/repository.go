package attendance

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, att *Attendance) error
	GetByUserAndDate(ctx context.Context, userID uint, date time.Time) (*Attendance, error)
	Update(ctx context.Context, att *Attendance) error
	GetByUser(ctx context.Context, userID uint, offset int) ([]Attendance, error)
	GetHistoryByUser(ctx context.Context, userID uint, from *time.Time, to *time.Time, limit int, offset int) ([]Attendance, error)
	GetByOrganisation(ctx context.Context, orgID uint, offset int) ([]Attendance, error)
}
