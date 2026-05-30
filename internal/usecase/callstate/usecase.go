package callstate

import (
	"context"
	domain "server/internal/domain/callstate"
)

type service struct{ repo domain.Repository }

func NewService(repo domain.Repository) domain.UseCase { return &service{repo: repo} }

func (s *service) CreateOrUpdate(ctx context.Context, state *domain.State) error {
	return s.repo.CreateOrUpdate(ctx, state)
}
func (s *service) Get(ctx context.Context, id uint) (*domain.State, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *service) CreateOffer(ctx context.Context, id uint, callingID, firstName, lastName, offer string) error {
	return s.repo.UpdateOffer(ctx, id, callingID, firstName, lastName, offer)
}
func (s *service) EndCall(ctx context.Context, id uint) error { return s.repo.EndCall(ctx, id) }
