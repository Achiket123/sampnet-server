package calendar

import (
	"context"
	"time"
)

type GetEventsRequest struct {
	OrgID     uint
	UserID    *uint
	StartDate time.Time
	EndDate   time.Time
}

type UseCase interface {
	GetEvents(ctx context.Context, req GetEventsRequest) ([]CalendarEvent, error)
}
