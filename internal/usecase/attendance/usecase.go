package attendance

import (
	"context"
	domain "server/internal/domain/attendance"
	"time"
	"errors"
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.UseCase {
	return &service{repo: repo}
}

func (s *service) RecordAttendance(ctx context.Context, att *domain.Attendance) error {
	if att.Date.IsZero() {
		att.Date = time.Now()
	}
	return s.repo.Create(ctx, att)
}

func (s *service) UpdateAttendance(ctx context.Context, userID uint, checkOutTime *time.Time, checkOutPhoto int) (*domain.Attendance, error) {
	today := time.Now()
	att, err := s.repo.GetByUserAndDate(ctx, userID, today)
	if err != nil {
		return nil, errors.New("PLEASE CHECK IN FIRST")
	}

	att.CheckOutTime = checkOutTime
	att.CheckOutPhoto = checkOutPhoto

	if err := s.repo.Update(ctx, att); err != nil {
		return nil, err
	}
	return att, nil
}

func (s *service) GetAttendanceByUser(ctx context.Context, userID uint, offset int) ([]domain.Attendance, error) {
	return s.repo.GetByUser(ctx, userID, offset)
}

func (s *service) GetAttendanceByOrganisation(ctx context.Context, orgID uint, offset int) ([]domain.Attendance, error) {
	return s.repo.GetByOrganisation(ctx, orgID, offset)
}

func (s *service) GetAttendanceByDateAndUser(ctx context.Context, userID uint, date time.Time) (*domain.Attendance, error) {
	return s.repo.GetByUserAndDate(ctx, userID, date)
}
