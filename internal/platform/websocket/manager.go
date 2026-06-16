package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Envelope is the standard JSON structure for all WebSocket messages
type Envelope struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Manager handles application-level WebSocket operations using the Hub
type Manager struct {
	Hub         *Hub
	RedisClient *redis.Client
	Channel     string
}

// NewManager creates a new WebSocketManager
func NewManager(hub *Hub, redisClient *redis.Client, channel string) *Manager {
	return &Manager{
		Hub:         hub,
		RedisClient: redisClient,
		Channel:     channel,
	}
}

type managerPubSubMessage struct {
	TargetUserID string `json:"target_user_id"`
	TargetRoomID string `json:"target_room_id"`
	Payload      []byte `json:"payload"`
}

// SendToUser marshals the payload into an envelope and sends it to the specific user via the Hub
func (m *Manager) SendToUser(userID string, eventType string, payload interface{}) error {
	envelope := Envelope{
		Type:    eventType,
		Payload: payload,
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal websocket envelope: %w", err)
	}

	pubMsg := managerPubSubMessage{
		TargetUserID: userID,
		Payload:      data,
	}

	msgData, err := json.Marshal(pubMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal pubsub message: %w", err)
	}

	// Publish to Redis channel with background context
	err = m.RedisClient.Publish(context.Background(), m.Channel, msgData).Err()
	if err != nil {
		log.Printf("Error publishing websocket message to Redis: %v", err)
		m.Hub.Broadcast <- &HubMessage{
			TargetUserID: userID,
			Payload:      data,
		}
	}

	return nil
}

// SendToRoom marshals the payload into an envelope and sends it to all clients in the room via the Hub
func (m *Manager) SendToRoom(roomID string, eventType string, payload interface{}) error {
	envelope := Envelope{
		Type:    eventType,
		Payload: payload,
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal websocket envelope: %w", err)
	}

	pubMsg := managerPubSubMessage{
		TargetRoomID: roomID,
		Payload:      data,
	}

	msgData, err := json.Marshal(pubMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal pubsub message: %w", err)
	}

	// Publish to Redis channel with background context
	err = m.RedisClient.Publish(context.Background(), m.Channel, msgData).Err()
	if err != nil {
		log.Printf("Error publishing room message to Redis: %v", err)
		m.Hub.Broadcast <- &HubMessage{
			TargetRoomID: roomID,
			Payload:      data,
		}
	}

	return nil
}

// IsOnline checks if a specific user is currently connected to the WebSocket
func (m *Manager) IsOnline(userID string) bool {
	return m.Hub.IsOnline(userID)
}

// GetOnlineUsers filters a list of user IDs to return only those currently online
func (m *Manager) GetOnlineUsers(userIDs []string) []string {
	online := make([]string, 0)
	for _, id := range userIDs {
		if m.IsOnline(id) {
			online = append(online, id)
		}
	}
	return online
}

