package calendar

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "server/internal/domain/calendar"
)

var (
	ErrStartDateAfterEndDate = errors.New("start date cannot be after end date")
	ErrDateRangeTooLarge     = errors.New("date range cannot exceed 3 months")
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.UseCase {
	return &service{repo: repo}
}

func (s *service) GetEvents(ctx context.Context, req domain.GetEventsRequest) ([]domain.CalendarEvent, error) {
	if req.StartDate.After(req.EndDate) {
		return nil, ErrStartDateAfterEndDate
	}

	// 95 days max (approx. 3 months)
	if req.EndDate.Sub(req.StartDate) > 95*24*time.Hour {
		return nil, ErrDateRangeTooLarge
	}

	events, err := s.repo.GetCalendarEvents(ctx, req.OrgID, req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("usecase failed: %w", err)
	}

	return events, nil
}