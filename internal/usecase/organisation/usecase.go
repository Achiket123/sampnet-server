package organisation

import (
	"context"
	"errors"
	domain "server/internal/domain/organisation"
)

var ErrInvalidID = errors.New("invalid organisation id")

// Service implements organisation application business rules.
type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, org *domain.Entity, ownerUserID uint) (*domain.OwnerEmployeeRow, error) {
	if ownerUserID == 0 {
		return nil, ErrInvalidID
	}
	org.BossID = ownerUserID
	
	return s.repo.CreateWithOwner(ctx, org, ownerUserID)
}

func (s *Service) Get(ctx context.Context, id uint) (*domain.Entity, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, org *domain.Entity) error {
	if org == nil || org.ID == 0 {
		return ErrInvalidID
	}
	return s.repo.Update(ctx, org)
}
