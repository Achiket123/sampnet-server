package messages

import "context"

type Repository interface {
	Create(ctx context.Context, msg *Message) error
	ListByRoom(ctx context.Context, roomID string) ([]Message, error)
	MarkSeen(ctx context.Context, roomID string, receiverID string) error
}
