package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/domain/notifications"

	"github.com/redis/go-redis/v9"
)

type RedisBroadcaster struct {
	hub         *Hub
	redisClient *redis.Client
	channel     string
}

func NewRedisBroadcaster(hub *Hub, redisClient *redis.Client, channel string) notifications.Broadcaster {
	return &RedisBroadcaster{
		hub:         hub,
		redisClient: redisClient,
		channel:     channel,
	}
}

func (b *RedisBroadcaster) BroadcastNotification(ctx context.Context, n *notifications.Notification) error {
	payload, err := json.Marshal(n)
	if err != nil {
		log.Printf("Error marshaling notification for pub/sub: %v", err)
		return err
	}

	// Publish message to Redis pub/sub channel
	err = b.redisClient.Publish(ctx, b.channel, payload).Err()
	if err != nil {
		log.Printf("Error publishing notification to Redis: %v", err)
		return err
	}

	return nil
}

func StartSubscriber(ctx context.Context, redisClient *redis.Client, channel string, hub *Hub) {
	pubsub := redisClient.Subscribe(ctx, channel)

	// Clean up pubsub subscription when context is cancelled
	go func() {
		<-ctx.Done()
		pubsub.Close()
	}()

	ch := pubsub.Channel()
	log.Printf("Started Redis Pub/Sub subscriber on channel: %s", channel)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping Redis Pub/Sub subscriber on channel %s due to context cancellation", channel)
			return
		case msg, ok := <-ch:
			if !ok {
				log.Printf("Redis Pub/Sub channel %s closed", channel)
				return
			}

			var n notifications.Notification
			if err := json.Unmarshal([]byte(msg.Payload), &n); err != nil {
				log.Printf("Error unmarshaling Redis pub/sub message: %v", err)
				continue
			}

			// Broadcast to the local WebSocket connections
			hub.Broadcast <- &HubMessage{
				TargetUserID: fmt.Sprintf("%d", n.UserID),
				Payload:      []byte(msg.Payload),
			}
		}
	}
}
