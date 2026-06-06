package websocket

import (
	"encoding/json"
	"fmt"
)

// Envelope is the standard JSON structure for all WebSocket messages
type Envelope struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Manager handles application-level WebSocket operations using the Hub
type Manager struct {
	Hub *Hub
}

// NewManager creates a new WebSocketManager
func NewManager(hub *Hub) *Manager {
	return &Manager{Hub: hub}
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

	m.Hub.Broadcast <- &HubMessage{
		TargetUserID: userID,
		Payload:      data,
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

	m.Hub.Broadcast <- &HubMessage{
		TargetRoomID: roomID,
		Payload:      data,
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
