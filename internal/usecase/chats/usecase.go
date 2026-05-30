package chats

import (
	"context"
	domain "server/internal/domain/chats"
)

type service struct{ repo domain.Repository }

func NewService(repo domain.Repository) domain.UseCase { return &service{repo: repo} }

func (s *service) CreateChat(ctx context.Context, chat *domain.Chat) error {
	return s.repo.Create(ctx, chat)
}
func (s *service) GetChats(ctx context.Context, organisationID uint) ([]domain.Chat, error) {
	return s.repo.ListByOrganisation(ctx, organisationID)
}
