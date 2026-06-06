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

func (s *service) GetUserChats(ctx context.Context, userID uint, orgID uint) ([]domain.Chat, error) {
	return s.repo.ListByUser(ctx, userID, orgID)
}

func (s *service) GetOrCreateDM(ctx context.Context, userA uint, userB uint, orgID uint) (*domain.Chat, error) {
	return s.repo.GetOrCreateDM(ctx, userA, userB, orgID)
}

func (s *service) GetChatByID(ctx context.Context, chatID uint) (*domain.Chat, error) {
	return s.repo.GetByID(ctx, chatID)
}
