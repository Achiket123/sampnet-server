package messages

import (
	"context"
	"errors"
	"strconv"
	"time"

	chatDomain "server/internal/domain/chats"
	domain "server/internal/domain/messages"
	ws "server/internal/platform/websocket"
)

var ErrInvalidUserID = errors.New("invalid user id")

type service struct {
	repo      domain.Repository
	chatRepo  chatDomain.Repository
	wsManager *ws.Manager
}

func NewService(repo domain.Repository, chatRepo chatDomain.Repository, wsManager *ws.Manager) domain.UseCase {
	return &service{repo: repo, chatRepo: chatRepo, wsManager: wsManager}
}

func (s *service) GetMessages(ctx context.Context, roomID string, cursor string, limit int) (*domain.CursorPage, error) {
	return s.repo.ListByRoomCursor(ctx, roomID, cursor, limit)
}

func (s *service) SendMessage(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now().UTC()
	}

	createdMsg, err := s.repo.Create(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Fetch the chat by room ID to update its metadata
	chat, err := s.chatRepo.GetByRoomID(ctx, createdMsg.RoomID)
	if err == nil {
		ts := createdMsg.CreatedAt.UTC().Format(time.RFC3339)
		_ = s.chatRepo.UpdateLastMessage(ctx, chat.ID, createdMsg.Message, &ts)
		_ = s.chatRepo.IncrementMessageCount(ctx, chat.ID)

		// Get chat participants to push notification
		participants, err := s.chatRepo.GetParticipants(ctx, chat.ID)
		if err == nil {
			for _, p := range participants {
				participantUserIDStr := strconv.FormatUint(uint64(p.UserID), 10)
				if participantUserIDStr != createdMsg.SenderID {
					_ = s.chatRepo.UpdateUnreadCount(ctx, chat.ID, p.UserID, 1)
				}
			}
			// Push via WebSocket to the room
			if s.wsManager != nil {
				_ = s.wsManager.SendToRoom(chat.RoomID, "new_message", createdMsg)
			}
		}
	}

	return createdMsg, nil
}

func (s *service) MarkSeen(ctx context.Context, roomID string, receiverID string) error {
	err := s.repo.MarkSeen(ctx, roomID, receiverID)
	if err == nil {
		// Reset unread count for this user in this chat
		chatID, parseErr := strconv.ParseUint(roomID, 10, 64)
		userID, parseErr2 := strconv.ParseUint(receiverID, 10, 64)
		if parseErr == nil && parseErr2 == nil {
			_ = s.chatRepo.ResetUnreadCount(ctx, uint(chatID), uint(userID))
		}
	}
	return err
}

func (s *service) DeleteMessage(ctx context.Context, messageID uint, requestingUserID string) error {
	return s.repo.DeleteMessage(ctx, messageID, requestingUserID)
}
