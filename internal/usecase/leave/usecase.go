package leave

import (
	"context"
	"fmt"
	"server/internal/domain/leave"
	"server/internal/domain/notifications"
	"time"
)

type service struct {
	repo           leave.Repository
	notificationUc notifications.UseCase
}

func NewService(repo leave.Repository, notificationUc notifications.UseCase) leave.UseCase {
	return &service{
		repo:           repo,
		notificationUc: notificationUc,
	}
}

func (s *service) RequestLeave(ctx context.Context, employeeID uint, orgID uint, leaveType string, startDate time.Time, endDate time.Time, reason string) (*leave.Leave, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	if startDate.Before(today) {
		return nil, fmt.Errorf("start date cannot be in the past")
	}
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end date must be on or after start date")
	}

	totalDays := int(endDate.Sub(startDate).Hours()/24) + 1

	l := &leave.Leave{
		EmployeeID:     employeeID,
		OrganisationID: orgID,
		LeaveType:      leaveType,
		StartDate:      startDate,
		EndDate:        endDate,
		TotalDays:      totalDays,
		Reason:         reason,
		Status:         "pending",
	}

	if err := s.repo.Create(ctx, l); err != nil {
		return nil, err
	}

	// Send confirmation to employee
	go func() {
		// TODO: Look up managers for this organisation and notify them as well.
		// Requires employee repository dependency.
		s.notificationUc.CreateNotification(
			context.Background(),
			employeeID,
			orgID,
			"Leave Request Submitted",
			fmt.Sprintf("Your %s leave request from %s to %s has been submitted.", leaveType, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
			"leave_request",
			"",
		)
	}()

	return l, nil
}

func (s *service) GetLeave(ctx context.Context, id uint) (*leave.Leave, error) {
	l, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if l == nil {
		return nil, leave.ErrLeaveNotFound
	}
	return l, nil
}

func (s *service) GetMyLeaves(ctx context.Context, employeeID uint, offset int) ([]leave.Leave, error) {
	return s.repo.GetByEmployee(ctx, employeeID, offset)
}

func (s *service) GetOrganisationLeaves(ctx context.Context, orgID uint, status string, offset int) ([]leave.Leave, error) {
	return s.repo.GetByOrganisation(ctx, orgID, status, offset)
}

func (s *service) ApproveLeave(ctx context.Context, leaveID uint, managerID uint, managerNote string) error {
	l, err := s.repo.GetByID(ctx, leaveID)
	if err != nil {
		return err
	}
	if l == nil {
		return leave.ErrLeaveNotFound
	}

	if l.Status != "pending" {
		return fmt.Errorf("leave is not in a pending state")
	}

	l.Status = "approved"
	l.ManagerID = managerID
	l.ManagerNote = managerNote

	if err := s.repo.Update(ctx, l); err != nil {
		return err
	}

	go s.notificationUc.CreateNotification(
		context.Background(),
		l.EmployeeID,
		l.OrganisationID,
		"Leave Approved",
		fmt.Sprintf("Your leave request for %s has been approved.", l.LeaveType),
		"leave_approved",
		"",
	)

	return nil
}

func (s *service) RejectLeave(ctx context.Context, leaveID uint, managerID uint, managerNote string) error {
	l, err := s.repo.GetByID(ctx, leaveID)
	if err != nil {
		return err
	}
	if l == nil {
		return leave.ErrLeaveNotFound
	}

	if l.Status != "pending" {
		return fmt.Errorf("leave is not in a pending state")
	}

	l.Status = "rejected"
	l.ManagerID = managerID
	l.ManagerNote = managerNote

	if err := s.repo.Update(ctx, l); err != nil {
		return err
	}

	go s.notificationUc.CreateNotification(
		context.Background(),
		l.EmployeeID,
		l.OrganisationID,
		"Leave Rejected",
		fmt.Sprintf("Your leave request for %s has been rejected.", l.LeaveType),
		"leave_rejected",
		"",
	)

	return nil
}

func (s *service) CancelLeave(ctx context.Context, leaveID uint, employeeID uint) error {
	l, err := s.repo.GetByID(ctx, leaveID)
	if err != nil {
		return err
	}
	if l == nil {
		return leave.ErrLeaveNotFound
	}

	if l.EmployeeID != employeeID {
		return leave.ErrNotAuthorized
	}

	if l.Status != "pending" {
		return fmt.Errorf("cannot cancel an already processed leave")
	}

	l.Status = "cancelled"
	return s.repo.Update(ctx, l)
}
