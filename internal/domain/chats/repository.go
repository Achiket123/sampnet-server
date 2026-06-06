package chats

import "context"

type Repository interface {
	Create(ctx context.Context, chat *Chat) error
	ListByOrganisation(ctx context.Context, organisationID uint) ([]Chat, error)
	ListByUser(ctx context.Context, userID uint, orgID uint) ([]Chat, error)
	GetByID(ctx context.Context, chatID uint) (*Chat, error)
	GetByRoomID(ctx context.Context, roomID string) (*Chat, error)
	UpdateLastMessage(ctx context.Context, id uint, message string, at *string) error
	IncrementMessageCount(ctx context.Context, id uint) error
	GetOrCreateDM(ctx context.Context, userA uint, userB uint, orgID uint) (*Chat, error)
	AddParticipant(ctx context.Context, chatID uint, userID uint) error
	UpdateUnreadCount(ctx context.Context, chatID uint, userID uint, delta int) error
	ResetUnreadCount(ctx context.Context, chatID uint, userID uint) error
	GetParticipants(ctx context.Context, chatID uint) ([]ChatParticipant, error)
}
