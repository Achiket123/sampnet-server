package models

import "time"

type Chat struct {
	ID                   uint       `gorm:"primaryKey;autoIncrement:false" json:"id"`
	FirstName            string     `gorm:"size:100;not null" json:"first_name"`
	LastName             string     `gorm:"size:100" json:"last_name"`
	Email                string     `gorm:"size:255" json:"email"`
	OrganisationID       uint       `gorm:"index;not null" json:"organisation_id"`
	LastMessage          string     `gorm:"type:text" json:"last_message"`
	LastMessageTimestamp *time.Time `json:"last_message_timestamp"`
	NumberOfMessage      int        `gorm:"default:0" json:"number_of_message"`
}

type ChatMessage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoomID     string    `gorm:"size:255;index;not null" json:"room_id"`
	SenderID   string    `gorm:"size:50;index;not null" json:"sender_id"`
	ReceiverID string    `gorm:"size:50;index;not null" json:"receiver_id"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	IsSeen     bool      `gorm:"default:false" json:"is_seen"`
	TimeStamp  time.Time `gorm:"index" json:"time_stamp"`
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
