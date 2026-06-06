package chats

import "time"

type Chat struct {
	ID             uint              `json:"id"`
	RoomID         string            `json:"room_id"`
	OrganisationID uint              `json:"organisation_id"`
	Name           string            `json:"name"`
	IsGroup        bool              `json:"is_group"`
	CreatedBy      uint              `json:"created_by"`
	LastMessage    string            `json:"last_message"`
	LastMessageAt  *time.Time        `json:"last_message_at"`
	MessageCount   int               `json:"message_count"`
	Participants   []ChatParticipant `json:"participants,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type ChatParticipant struct {
	ChatID            uint       `json:"chat_id"`
	UserID            uint       `json:"user_id"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	AvatarURL         string     `json:"avatar_url"`
	UnreadCount       int        `json:"unread_count"`
	LastReadMessageID *uint      `json:"last_read_message_id"`
	JoinedAt          time.Time  `json:"joined_at"`
}
