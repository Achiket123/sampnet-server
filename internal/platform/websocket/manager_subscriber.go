package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

// StartManagerSubscriber listens to the Redis channel for WebSocket message envelopes to broadcast locally.
func StartManagerSubscriber(ctx context.Context, redisClient *redis.Client, channel string, hub *Hub) {
	pubsub := redisClient.Subscribe(ctx, channel)

	// Close the subscription when the context is cancelled.
	go func() {
		<-ctx.Done()
		pubsub.Close()
	}()

	ch := pubsub.Channel()
	log.Printf("Started WebSocket Manager subscriber on channel: %s", channel)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping WebSocket Manager subscriber on channel %s due to context cancellation", channel)
			return
		case msg, ok := <-ch:
			if !ok {
				log.Printf("WebSocket Manager channel %s closed", channel)
				return
			}

			var payload managerPubSubMessage
			if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
				log.Printf("Error unmarshaling WebSocket Manager pub/sub message: %v", err)
				continue
			}

			hub.Broadcast <- &HubMessage{
				TargetUserID: payload.TargetUserID,
				TargetRoomID: payload.TargetRoomID,
				Payload:      payload.Payload,
			}
		}
	}
}
