package teams

import (
	"context"
	domain "server/internal/domain/teams"
)

type service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) domain.UseCase {
	return &service{repo: repo}
}

func (s *service) CreateTeam(ctx context.Context, team *domain.Team, members []uint) error {
	team.IsActive = true
	if err := s.repo.Create(ctx, team); err != nil {
		return err
	}

	for _, memberID := range members {
		var role string
		if team.TeamLead == memberID {
			role = "Team Lead"
		} else {
			role = "Member"
		}
		member := &domain.TeamMember{
			UserID:   memberID,
			TeamID:   team.ID,
			Role:     role,
			IsActive: true,
			IsLeader: team.TeamLead == memberID,
		}
		if err := s.repo.CreateMember(ctx, member); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetTeam(ctx context.Context, id uint) (*domain.Team, []domain.TeamMember, error) {
	team, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	members, err := s.repo.GetMembersByTeam(ctx, id)
	if err != nil {
		return team, nil, err
	}
	return team, members, nil
}

func (s *service) UpdateTeam(ctx context.Context, team *domain.Team) error {
	return s.repo.Update(ctx, team)
}

func (s *service) DeleteTeam(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetTeamsByOrganisation(ctx context.Context, orgID uint) ([]domain.Team, error) {
	return s.repo.GetByOrganisation(ctx, orgID)
}

func (s *service) GetTeams(ctx context.Context) ([]domain.Team, error) {
	return s.repo.GetAll(ctx)
}
