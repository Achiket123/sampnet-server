package messages

import "context"

type UseCase interface {
	GetMessages(ctx context.Context, roomID string, cursor string, limit int) (*CursorPage, error)
	SendMessage(ctx context.Context, msg *Message) (*Message, error)
	MarkSeen(ctx context.Context, roomID string, receiverID string) error
	DeleteMessage(ctx context.Context, messageID uint, requestingUserID string) error
}
