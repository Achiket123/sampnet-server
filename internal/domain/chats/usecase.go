package chats

import "context"

type UseCase interface {
	CreateChat(ctx context.Context, chat *Chat) error
	GetChats(ctx context.Context, organisationID uint) ([]Chat, error)
}
