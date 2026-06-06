package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/domain/notifications"
)

type Broadcaster struct {
	hub *Hub
}

func NewBroadcaster(hub *Hub) notifications.Broadcaster {
	return &Broadcaster{hub: hub}
}

func (b *Broadcaster) BroadcastNotification(ctx context.Context, n *notifications.Notification) error {
	message, err := json.Marshal(n)
	if err != nil {
		log.Printf("Error marshaling notification: %v", err)
		return err
	}

	// Send to the specific user's websocket connections
	b.hub.Broadcast <- &HubMessage{
		TargetUserID: fmt.Sprintf("%d", n.UserID),
		Payload:      message,
	}
	return nil
}