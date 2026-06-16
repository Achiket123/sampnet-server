package invites

import (
	"context"
	authDomain "server/internal/domain/auth"
)

type UseCase interface {
	InviteEmployee(ctx context.Context, invite *EmployeeInvite) error
	AcceptInvite(ctx context.Context, token string, password string) (authDomain.TokenPair, error)
	GetInvitesByOrg(ctx context.Context, orgID uint) ([]EmployeeInvite, error)
	ResendInvite(ctx context.Context, email string, orgID uint) error
}
