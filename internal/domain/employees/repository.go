package employees

import "context"

type Repository interface {
	GetByUserID(ctx context.Context, userID uint) (*Employee, error)
	Create(ctx context.Context, employee *Employee) error
	GetEmployeesByOrg(ctx context.Context, orgID uint) ([]Employee, error)
	GetByID(ctx context.Context, id uint) (*Employee, error)
	Update(ctx context.Context, employee *Employee) error
	Delete(ctx context.Context, id uint) error
	Search(ctx context.Context, query string) ([]Employee, error)

	CreateManager(ctx context.Context, manager *Manager) error
	GetManagerByUserID(ctx context.Context, userID uint) (*Manager, error)

	CreateBoss(ctx context.Context, boss *Boss) error
	GetBossByUserID(ctx context.Context, userID uint) (*Boss, error)
}
