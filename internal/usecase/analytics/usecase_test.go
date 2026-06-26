package analytics

import (
	"context"
	"errors"
	"testing"

	domain "server/internal/domain/analytics"
	empDomain "server/internal/domain/employees"
)

type mockRepository struct {
	domain.Repository
	getEmployeeAnalyticsFunc  func(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error)
	getOrgEmployeeMonitorFunc func(ctx context.Context, orgID uint) ([]domain.EmployeeMonitorResponse, error)
}

func (m *mockRepository) GetEmployeeAnalytics(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error) {
	if m.getEmployeeAnalyticsFunc != nil {
		return m.getEmployeeAnalyticsFunc(ctx, userID, orgID, period)
	}
	return nil, nil
}

func (m *mockRepository) GetOrgEmployeeMonitor(ctx context.Context, orgID uint) ([]domain.EmployeeMonitorResponse, error) {
	if m.getOrgEmployeeMonitorFunc != nil {
		return m.getOrgEmployeeMonitorFunc(ctx, orgID)
	}
	return nil, nil
}

func TestGetEmployeeAnalytics_DelegatesToRepository(t *testing.T) {
	expected := &domain.EmployeeAnalyticsSummary{
		EmployeeInfo: empDomain.Employee{UserID: 11},
	}

	var receivedUserID uint
	var receivedOrgID uint
	var receivedPeriod string

	repo := &mockRepository{
		getEmployeeAnalyticsFunc: func(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error) {
			receivedUserID = userID
			receivedOrgID = orgID
			receivedPeriod = period
			return expected, nil
		},
	}

	svc := NewService(repo)
	got, err := svc.GetEmployeeAnalytics(context.Background(), 11, 22, "week")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != expected {
		t.Fatalf("expected returned summary to be passed through unchanged")
	}
	if receivedUserID != 11 || receivedOrgID != 22 || receivedPeriod != "week" {
		t.Fatalf("unexpected repository arguments: userID=%d orgID=%d period=%q", receivedUserID, receivedOrgID, receivedPeriod)
	}
}

func TestGetEmployeeAnalytics_PropagatesRepositoryError(t *testing.T) {
	repoErr := errors.New("repository failed")
	repo := &mockRepository{
		getEmployeeAnalyticsFunc: func(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error) {
			return nil, repoErr
		},
	}

	svc := NewService(repo)
	_, err := svc.GetEmployeeAnalytics(context.Background(), 11, 22, "month")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}

func TestGetOrgEmployeeMonitor_DelegatesToRepository(t *testing.T) {
	expected := []domain.EmployeeMonitorResponse{
		{EmployeeInfo: empDomain.Employee{UserID: 44}},
	}

	var receivedOrgID uint
	repo := &mockRepository{
		getOrgEmployeeMonitorFunc: func(ctx context.Context, orgID uint) ([]domain.EmployeeMonitorResponse, error) {
			receivedOrgID = orgID
			return expected, nil
		},
	}

	svc := NewService(repo)
	got, err := svc.GetOrgEmployeeMonitor(context.Background(), 77)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].EmployeeInfo.UserID != 44 {
		t.Fatalf("unexpected returned monitor data: %#v", got)
	}
	if receivedOrgID != 77 {
		t.Fatalf("expected orgID 77, got %d", receivedOrgID)
	}
}
