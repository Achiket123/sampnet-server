package employees

import "context"

type UseCase interface {
	AddEmployee(ctx context.Context, emp *Employee) error
	GetEmployees(ctx context.Context, orgID uint) ([]Employee, error)
	GetEmployee(ctx context.Context, id uint) (*Employee, error)
	UpdateEmployee(ctx context.Context, emp *Employee) error
	DeleteEmployee(ctx context.Context, id uint) error
	SearchEmployees(ctx context.Context, query string) ([]Employee, error)
	
	MakeManager(ctx context.Context, manager *Manager) error
	IsEmployeeOrManager(ctx context.Context, userID uint) (string, interface{}, string, error)

	CreateBoss(ctx context.Context, boss *Boss) error
	GetBoss(ctx context.Context, userID uint) (*Boss, error)
}
