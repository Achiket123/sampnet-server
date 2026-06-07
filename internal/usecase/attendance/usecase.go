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

func (s *service) GetEmployeeAttendanceHistory(ctx context.Context, userID uint, from *time.Time, to *time.Time, limit int, offset int) (*domain.AttendanceHistory, error) {
	records, err := s.repo.GetHistoryByUser(ctx, userID, from, to, limit, offset)
	if err != nil {
		return nil, err
	}

	totalDaysPresent := len(records)
	var totalDaysAbsent int
	
	if from != nil && to != nil {
		days := int(to.Sub(*from).Hours() / 24) + 1
		totalDaysAbsent = days - totalDaysPresent
		if totalDaysAbsent < 0 {
			totalDaysAbsent = 0
		}
	}

	var totalDuration int
	var onTimeCount int
	var totalCheckInMinutes int
	var checkInCount int

	for _, rec := range records {
		if rec.DurationMinutes != nil {
			totalDuration += *rec.DurationMinutes
		}
		if rec.CheckInTime != nil {
			checkInCount++
			totalCheckInMinutes += rec.CheckInTime.Hour()*60 + rec.CheckInTime.Minute()
			
			// Assuming on-time is before 09:30 AM
			if rec.CheckInTime.Hour() < 9 || (rec.CheckInTime.Hour() == 9 && rec.CheckInTime.Minute() <= 30) {
				onTimeCount++
			}
		}
	}

	var averageDuration int
	if totalDaysPresent > 0 {
		averageDuration = totalDuration / totalDaysPresent
	}

	var avgCheckInStr string
	if checkInCount > 0 {
		avgMins := totalCheckInMinutes / checkInCount
		avgCheckInTime := time.Date(0, 1, 1, avgMins/60, avgMins%60, 0, 0, time.UTC)
		avgCheckInStr = avgCheckInTime.Format("03:04 PM")
	}

	var onTimePercentage float64
	if totalDaysPresent > 0 {
		onTimePercentage = float64(onTimeCount) / float64(totalDaysPresent) * 100
	}

	summary := domain.AttendanceSummary{
		TotalDaysPresent:       totalDaysPresent,
		TotalDaysAbsent:        totalDaysAbsent,
		AverageCheckInTime:     avgCheckInStr,
		AverageDurationMinutes: averageDuration,
		OnTimePercentage:       onTimePercentage,
	}

	return &domain.AttendanceHistory{
		Records: records,
		Summary: summary,
	}, nil
}
