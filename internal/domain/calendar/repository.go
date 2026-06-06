package calendar

import (
	"context"
	"time"
)

type Repository interface {
	GetCalendarEvents(ctx context.Context, orgID uint, userID *uint, startDate time.Time, endDate time.Time) ([]CalendarEvent, error)
}