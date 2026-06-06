package chats

import "context"

type UseCase interface {
	CreateChat(ctx context.Context, chat *Chat) error
	GetChats(ctx context.Context, organisationID uint) ([]Chat, error)
	GetUserChats(ctx context.Context, userID uint, orgID uint) ([]Chat, error)
	GetOrCreateDM(ctx context.Context, userA uint, userB uint, orgID uint) (*Chat, error)
	GetChatByID(ctx context.Context, chatID uint) (*Chat, error)
}
