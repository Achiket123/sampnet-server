package leave

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, leave *Leave) error
	GetByID(ctx context.Context, id uint) (*Leave, error)
	GetByEmployee(ctx context.Context, employeeID uint, offset int) ([]Leave, error)
	GetByOrganisation(ctx context.Context, orgID uint, status string, offset int) ([]Leave, error)
	Update(ctx context.Context, leave *Leave) error
	Delete(ctx context.Context, id uint) error
}
