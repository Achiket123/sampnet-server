package teams

import "context"

type UseCase interface {
	CreateTeam(ctx context.Context, team *Team, members []uint) error
	GetTeam(ctx context.Context, id uint) (*Team, []TeamMember, error)
	UpdateTeam(ctx context.Context, team *Team) error
	DeleteTeam(ctx context.Context, id uint) error
	GetTeamsByOrganisation(ctx context.Context, orgID uint) ([]Team, error)
	GetTeams(ctx context.Context) ([]Team, error)
}
