package notifications

import "time"

type Notification struct {
	ID               uint      `json:"id"`
	UserID           uint      `json:"user_id"`
	OrganisationID   uint      `json:"organisation_id"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`
	NotificationType string    `json:"notification_type"`
	Link             string    `json:"link"`
	CreatedAt        time.Time `json:"created_at"`
}
