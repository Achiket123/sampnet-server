package messages

import "context"

type UseCase interface {
	GetMessages(ctx context.Context, me uint, peerID string) ([]Message, error)
	SendMessage(ctx context.Context, msg *Message) error
}
