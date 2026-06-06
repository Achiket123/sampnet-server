package leave

import (
	"context"
	"server/internal/domain/leave"
	"server/internal/domain/notifications"
	"testing"
	"time"
)

type mockRepository struct {
	leave.Repository
	createFunc func(ctx context.Context, l *leave.Leave) error
	getByIDFunc func(ctx context.Context, id uint) (*leave.Leave, error)
	updateFunc func(ctx context.Context, l *leave.Leave) error
}

func (m *mockRepository) Create(ctx context.Context, l *leave.Leave) error {
	return m.createFunc(ctx, l)
}
func (m *mockRepository) GetByID(ctx context.Context, id uint) (*leave.Leave, error) {
	return m.getByIDFunc(ctx, id)
}
func (m *mockRepository) Update(ctx context.Context, l *leave.Leave) error {
	return m.updateFunc(ctx, l)
}

type mockNotificationUseCase struct {
	createFunc func(ctx context.Context, userID uint, organisationID uint, title string, message string, notificationType string, link string) error
}

func (m *mockNotificationUseCase) CreateNotification(ctx context.Context, userID uint, organisationID uint, title string, message string, notificationType string, link string) error {
	return m.createFunc(ctx, userID, organisationID, title, message, notificationType, link)
}

func (m *mockNotificationUseCase) GetNotifications(ctx context.Context, userID uint, offset int) ([]notifications.Notification, error) {
	return nil, nil
}
func (m *mockNotificationUseCase) MarkNotificationRead(ctx context.Context, notificationID uint) error {
	return nil
}
func (m *mockNotificationUseCase) MarkAllNotificationsRead(ctx context.Context, userID uint) error {
	return nil
}

func TestRequestLeave(t *testing.T) {
	repo := &mockRepository{
		createFunc: func(ctx context.Context, l *leave.Leave) error {
			l.ID = 1
			return nil
		},
	}
	notif := &mockNotificationUseCase{
		createFunc: func(ctx context.Context, userID uint, organisationID uint, title string, message string, notificationType string, link string) error {
			return nil
		},
	}

	service := NewService(repo, notif)

	startDate := time.Now().Add(24 * time.Hour)
	endDate := time.Now().Add(48 * time.Hour)

	l, err := service.RequestLeave(context.Background(), 1, 1, "annual", startDate, endDate, "Vacation")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if l.ID != 1 {
		t.Errorf("expected ID 1, got %d", l.ID)
	}

	if l.TotalDays != 2 {
		t.Errorf("expected 2 total days, got %d", l.TotalDays)
	}
}

func TestApproveLeave(t *testing.T) {
	tLeave := &leave.Leave{
		ID: 1,
		Status: "pending",
		EmployeeID: 1,
		OrganisationID: 1,
	}

	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id uint) (*leave.Leave, error) {
			return tLeave, nil
		},
		updateFunc: func(ctx context.Context, l *leave.Leave) error {
			return nil
		},
	}
	notif := &mockNotificationUseCase{
		createFunc: func(ctx context.Context, userID uint, organisationID uint, title string, message string, notificationType string, link string) error {
			return nil
		},
	}

	service := NewService(repo, notif)

	err := service.ApproveLeave(context.Background(), 1, 2, "Enjoy!")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tLeave.Status != "approved" {
		t.Errorf("expected status approved, got %s", tLeave.Status)
	}
}

func TestApproveLeave_InvalidState(t *testing.T) {
	tLeave := &leave.Leave{
		ID: 1,
		Status: "rejected",
	}

	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id uint) (*leave.Leave, error) {
			return tLeave, nil
		},
	}

	service := NewService(repo, nil)

	err := service.ApproveLeave(context.Background(), 1, 2, "")

	if err == nil {
		t.Fatal("expected error for non-pending leave, got nil")
	}
}

func TestCancelLeave(t *testing.T) {
	tLeave := &leave.Leave{
		ID: 1,
		Status: "pending",
		EmployeeID: 1,
	}

	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id uint) (*leave.Leave, error) {
			return tLeave, nil
		},
		updateFunc: func(ctx context.Context, l *leave.Leave) error {
			return nil
		},
	}

	service := NewService(repo, nil)

	err := service.CancelLeave(context.Background(), 1, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tLeave.Status != "cancelled" {
		t.Errorf("expected status cancelled, got %s", tLeave.Status)
	}
}

func TestCancelLeave_Unauthorized(t *testing.T) {
	tLeave := &leave.Leave{
		ID: 1,
		Status: "pending",
		EmployeeID: 1,
	}

	repo := &mockRepository{
		getByIDFunc: func(ctx context.Context, id uint) (*leave.Leave, error) {
			return tLeave, nil
		},
	}

	service := NewService(repo, nil)

	err := service.CancelLeave(context.Background(), 1, 99) // Wrong user

	if err != leave.ErrNotAuthorized {
		t.Fatalf("expected ErrNotAuthorized, got %v", err)
	}
}
