package messages

import "time"

type Message struct {
	ID         uint      `json:"id"`
	RoomID     string    `json:"room_id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Message    string    `json:"message"`
	IsSeen     bool      `json:"is_seen"`
	TimeStamp  time.Time `json:"time_stamp"`
}
