package messages

import "context"

type Repository interface {
	// Create inserts a new message and returns it with server-assigned ID and timestamp
	Create(ctx context.Context, msg *Message) (*Message, error)

	// ListByRoomCursor returns messages for a room using cursor-based pagination.
	// cursor is empty for the first page. Subsequent pages pass NextCursor.
	ListByRoomCursor(ctx context.Context, roomID string, cursor string, limit int) (*CursorPage, error)

	// MarkSeen marks all messages in roomID as seen where receiver_id = receiverID
	MarkSeen(ctx context.Context, roomID string, receiverID string) error

	// MarkSeenUpTo marks messages as seen up to and including a specific message ID
	MarkSeenUpTo(ctx context.Context, roomID string, receiverID string, upToMessageID uint) error

	// DeleteMessage soft-deletes a message
	DeleteMessage(ctx context.Context, messageID uint, requestingUserID string) error

	// GetByID fetches a single message by ID
	GetByID(ctx context.Context, id uint) (*Message, error)
}
