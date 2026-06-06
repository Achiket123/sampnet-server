package messages

import "time"

type Message struct {
	ID             uint       `json:"id"`
	RoomID         string     `json:"room_id"`
	SenderID       string     `json:"sender_id"`
	ReceiverID     string     `json:"receiver_id"`
	OrganisationID uint       `json:"organisation_id"`
	Message        string     `json:"message"`
	MessageType    string     `json:"message_type"`
	FileURL        *string    `json:"file_url"`
	FileName       *string    `json:"file_name"`
	FileSize       *int64     `json:"file_size"`
	IsSeen         bool       `json:"is_seen"`
	IsDeleted      bool       `json:"is_deleted"`
	ReplyToID      *uint      `json:"reply_to_id"`
	ReplyToMessage  *Message   `json:"reply_to_message,omitempty"`
	SenderFirstName string     `json:"sender_first_name"`
	SenderLastName  string     `json:"sender_last_name"`
	SenderAvatarURL *string    `json:"sender_avatar_url"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CursorPage struct {
	Messages   []Message `json:"messages"`
	NextCursor string    `json:"next_cursor"`
	HasMore    bool      `json:"has_more"`
}
