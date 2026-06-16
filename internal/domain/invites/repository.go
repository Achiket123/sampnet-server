package invites

import (
	"context"
	authDomain "server/internal/domain/auth"
	employeesDomain "server/internal/domain/employees"
)

type Repository interface {
	Create(ctx context.Context, invite *EmployeeInvite) error
	GetByToken(ctx context.Context, token string) (*EmployeeInvite, error)
	GetByEmail(ctx context.Context, email string, orgID uint) (*EmployeeInvite, error)
	MarkAccepted(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*EmployeeInvite, error)
	Update(ctx context.Context, invite *EmployeeInvite) error
	GetByOrg(ctx context.Context, orgID uint) ([]EmployeeInvite, error)
	AcceptInvite(ctx context.Context, inviteID uint, user *authDomain.User, employee *employeesDomain.Employee) error
}