package messages

import (
	"context"
	"errors"
	chatDomain "server/internal/domain/chats"
	domain "server/internal/domain/messages"
	"sort"
	"strconv"
	"time"
)

var ErrInvalidUserID = errors.New("invalid user id")

type service struct {
	repo     domain.Repository
	chatRepo chatDomain.Repository
}

func NewService(repo domain.Repository, chatRepo chatDomain.Repository) domain.UseCase {
	return &service{repo: repo, chatRepo: chatRepo}
}

func roomID(a, b string) string {
	ids := []string{a, b}
	sort.Strings(ids)
	return ids[0] + "-" + ids[1]
}

func (s *service) GetMessages(ctx context.Context, me uint, peerID string) ([]domain.Message, error) {
	mine := strconv.FormatUint(uint64(me), 10)
	r := roomID(peerID, mine)
	_ = s.repo.MarkSeen(ctx, r, mine)
	return s.repo.ListByRoom(ctx, r)
}

func (s *service) SendMessage(ctx context.Context, msg *domain.Message) error {
	receiver, err := strconv.ParseUint(msg.ReceiverID, 10, 64)
	if err != nil {
		return ErrInvalidUserID
	}

	msg.RoomID = roomID(msg.ReceiverID, msg.SenderID)
	if msg.TimeStamp.IsZero() {
		msg.TimeStamp = time.Now().UTC()
	}

	if err := s.chatRepo.IncrementMessageCount(ctx, uint(receiver)); err != nil {
		return err
	}
	ts := msg.TimeStamp.UTC().Format(time.RFC3339)
	if err := s.chatRepo.UpdateLastMessage(ctx, uint(receiver), msg.Message, &ts); err != nil {
		return err
	}
	return s.repo.Create(ctx, msg)
}
