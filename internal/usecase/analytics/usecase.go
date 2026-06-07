package analytics

import (
	"context"
	"server/internal/domain/analytics"
)

type service struct {
	repo analytics.Repository
}

func NewService(repo analytics.Repository) analytics.UseCase {
	return &service{repo: repo}
}

func (s *service) GetEmployeeAnalytics(ctx context.Context, userID uint, orgID uint, period string) (*analytics.EmployeeAnalyticsSummary, error) {
	return s.repo.GetEmployeeAnalytics(ctx, userID, orgID, period)
}

func (s *service) GetOrgEmployeeMonitor(ctx context.Context, orgID uint) ([]analytics.EmployeeMonitorResponse, error) {
	return s.repo.GetOrgEmployeeMonitor(ctx, orgID)
}