package teams

import "context"

type Repository interface {
	Create(ctx context.Context, team *Team) error
	GetByID(ctx context.Context, id uint) (*Team, error)
	Update(ctx context.Context, team *Team) error
	Delete(ctx context.Context, id uint) error
	GetByOrganisation(ctx context.Context, orgID uint) ([]Team, error)
	GetAll(ctx context.Context) ([]Team, error)

	CreateMember(ctx context.Context, member *TeamMember) error
	GetMembersByTeam(ctx context.Context, teamID uint) ([]TeamMember, error)
}
