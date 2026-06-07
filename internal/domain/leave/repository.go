package leave

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, leave *Leave) error
	GetByID(ctx context.Context, id uint) (*Leave, error)
	GetByEmployee(ctx context.Context, employeeID uint, offset int) ([]Leave, error)
	GetByOrganisation(ctx context.Context, orgID uint, status string, offset int) ([]Leave, error)
	GetHistoryByEmployee(ctx context.Context, employeeID uint, status string, from *time.Time, to *time.Time, limit int, offset int) ([]Leave, error)
	Update(ctx context.Context, leave *Leave) error
	Delete(ctx context.Context, id uint) error
}
