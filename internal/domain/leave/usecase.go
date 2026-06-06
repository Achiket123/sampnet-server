package leave

import (
	"context"
	"time"
)

type UseCase interface {
	RequestLeave(ctx context.Context, employeeID uint, orgID uint, leaveType string, startDate time.Time, endDate time.Time, reason string) (*Leave, error)
	GetLeave(ctx context.Context, id uint) (*Leave, error)
	GetMyLeaves(ctx context.Context, employeeID uint, offset int) ([]Leave, error)
	GetOrganisationLeaves(ctx context.Context, orgID uint, status string, offset int) ([]Leave, error)
	ApproveLeave(ctx context.Context, leaveID uint, managerID uint, managerNote string) error
	RejectLeave(ctx context.Context, leaveID uint, managerID uint, managerNote string) error
	CancelLeave(ctx context.Context, leaveID uint, employeeID uint) error
}
