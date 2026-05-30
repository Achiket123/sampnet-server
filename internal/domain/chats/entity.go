package chats

import "time"

type Chat struct {
	ID                   uint       `json:"id"`
	FirstName            string     `json:"first_name"`
	LastName             string     `json:"last_name"`
	Email                string     `json:"email"`
	OrganisationID       uint       `json:"organisation_id"`
	LastMessage          string     `json:"last_message"`
	LastMessageTimestamp *time.Time `json:"last_message_timestamp"`
	NumberOfMessage      int        `json:"number_of_message"`
}
