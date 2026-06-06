package models

import "time"

type Chat struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	RoomID         string     `gorm:"size:255;unique;not null" json:"room_id"`
	OrganisationID uint       `gorm:"index;not null" json:"organisation_id"`
	Name           string     `gorm:"size:500" json:"name"`
	IsGroup        bool       `gorm:"default:false" json:"is_group"`
	CreatedBy      *uint      `json:"created_by"`
	LastMessage    string     `gorm:"type:text" json:"last_message"`
	LastMessageAt  *time.Time `json:"last_message_at"`
	MessageCount   int        `gorm:"default:0" json:"message_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ChatParticipant struct {
	ChatID            uint      `gorm:"primaryKey" json:"chat_id"`
	UserID            uint      `gorm:"primaryKey" json:"user_id"`
	UnreadCount       int       `gorm:"default:0" json:"unread_count"`
	LastReadMessageID *uint     `json:"last_read_message_id"`
	JoinedAt          time.Time `gorm:"default:now()" json:"joined_at"`
}

type ChatMessage struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	RoomID         string     `gorm:"size:255;index;not null" json:"room_id"`
	SenderID       string     `gorm:"size:255;not null" json:"sender_id"`
	ReceiverID     string     `gorm:"size:255;not null" json:"receiver_id"`
	OrganisationID uint       `gorm:"index;not null" json:"organisation_id"`
	Message        string     `gorm:"type:text;not null" json:"message"`
	MessageType    string     `gorm:"size:50;not null;default:'text'" json:"message_type"`
	FileURL        *string    `gorm:"type:text" json:"file_url"`
	FileName       *string    `gorm:"size:500" json:"file_name"`
	FileSize       *int64     `json:"file_size"`
	IsSeen         bool       `gorm:"default:false" json:"is_seen"`
	IsDeleted      bool       `gorm:"default:false" json:"is_deleted"`
	ReplyToID      *uint      `json:"reply_to_id"`
	CreatedAt      time.Time  `gorm:"index;default:now()" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"default:now()" json:"updated_at"`
}


type CallState struct {
	ID               uint       `gorm:"primaryKey;autoIncrement:false" json:"id"`
	FirstName        string     `gorm:"size:100" json:"first_name"`
	LastName         string     `gorm:"size:100" json:"last_name"`
	Email            string     `gorm:"size:255" json:"email"`
	OrganisationID   uint       `gorm:"index;not null" json:"organisation_id"`
	InCall           bool       `gorm:"default:false" json:"in_call"`
	LastCall         *time.Time `json:"last_call"`
	CallingID        *string    `gorm:"size:50" json:"calling_id"`
	CallingFirstName *string    `gorm:"size:100" json:"calling_first_name"`
	CallingLastName  *string    `gorm:"size:100" json:"calling_last_name"`
	Offer            *string    `gorm:"type:text" json:"offer"`
}
