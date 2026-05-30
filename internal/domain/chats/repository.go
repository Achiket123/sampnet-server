package chats

import "context"

type Repository interface {
	Create(ctx context.Context, chat *Chat) error
	ListByOrganisation(ctx context.Context, organisationID uint) ([]Chat, error)
	UpdateLastMessage(ctx context.Context, id uint, message string, at *string) error
	IncrementMessageCount(ctx context.Context, id uint) error
}
