package analytics

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domain "server/internal/domain/analytics"
	empDomain "server/internal/domain/employees"

	"github.com/gin-gonic/gin"
)

type mockUseCase struct {
	domain.UseCase
	getEmployeeAnalyticsFunc  func(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error)
	getOrgEmployeeMonitorFunc func(ctx context.Context, orgID uint) ([]domain.EmployeeMonitorResponse, error)
}

func (m *mockUseCase) GetEmployeeAnalytics(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error) {
	if m.getEmployeeAnalyticsFunc != nil {
		return m.getEmployeeAnalyticsFunc(ctx, userID, orgID, period)
	}
	return nil, nil
}

func (m *mockUseCase) GetOrgEmployeeMonitor(ctx context.Context, orgID uint) ([]domain.EmployeeMonitorResponse, error) {
	if m.getOrgEmployeeMonitorFunc != nil {
		return m.getOrgEmployeeMonitorFunc(ctx, orgID)
	}
	return nil, nil
}

func TestGetEmployeeAnalytics_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(&mockUseCase{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/employee/abc?organisation_id=1", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "userId", Value: "abc"}}

	h.GetEmployeeAnalytics(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Invalid user ID") {
		t.Fatalf("expected invalid user id response, got %s", w.Body.String())
	}
}

func TestGetEmployeeAnalytics_DefaultPeriodAndSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var receivedUserID uint
	var receivedOrgID uint
	var receivedPeriod string

	expected := &domain.EmployeeAnalyticsSummary{
		EmployeeInfo: empDomain.Employee{UserID: 9},
	}
	h := NewHandler(&mockUseCase{
		getEmployeeAnalyticsFunc: func(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error) {
			receivedUserID = userID
			receivedOrgID = orgID
			receivedPeriod = period
			return expected, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/employee/9?organisation_id=12", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "userId", Value: "9"}}

	h.GetEmployeeAnalytics(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if receivedUserID != 9 || receivedOrgID != 12 || receivedPeriod != "month" {
		t.Fatalf("unexpected arguments: userID=%d orgID=%d period=%q", receivedUserID, receivedOrgID, receivedPeriod)
	}
	if !strings.Contains(w.Body.String(), `"user_id":9`) {
		t.Fatalf("expected summary payload in response, got %s", w.Body.String())
	}
}

func TestGetEmployeeAnalytics_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repoErr := errors.New("boom")
	h := NewHandler(&mockUseCase{
		getEmployeeAnalyticsFunc: func(ctx context.Context, userID uint, orgID uint, period string) (*domain.EmployeeAnalyticsSummary, error) {
			return nil, repoErr
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/employee/9?organisation_id=12&period=week", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "userId", Value: "9"}}

	h.GetEmployeeAnalytics(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
	if !strings.Contains(w.Body.String(), "Failed to fetch analytics") {
		t.Fatalf("expected failure response, got %s", w.Body.String())
	}
}

func TestGetOrgEmployeeMonitor_InvalidOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(&mockUseCase{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/organisation/abc/employees", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "orgId", Value: "abc"}}

	h.GetOrgEmployeeMonitor(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetOrgEmployeeMonitor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expected := []domain.EmployeeMonitorResponse{
		{EmployeeInfo: empDomain.Employee{UserID: 4}},
	}
	h := NewHandler(&mockUseCase{
		getOrgEmployeeMonitorFunc: func(ctx context.Context, orgID uint) ([]domain.EmployeeMonitorResponse, error) {
			if orgID != 42 {
				t.Fatalf("expected orgID 42, got %d", orgID)
			}
			return expected, nil
		},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/organisation/42/employees", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "orgId", Value: "42"}}

	h.GetOrgEmployeeMonitor(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), `"user_id":4`) {
		t.Fatalf("expected monitor payload in response, got %s", w.Body.String())
	}
}
