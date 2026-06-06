package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// PresenceMessage is sent to clients when the list of online users changes
type PresenceMessage struct {
	OnlineUsers []uint `json:"online_users"`
}

type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Mutex for safe concurrent map access.
	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: user %d, employee %d, org %d", client.UserID, client.EmployeeID, client.OrganisationID)
			h.broadcastPresence(client.OrganisationID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Client unregistered: user %d, employee %d, org %d", client.UserID, client.EmployeeID, client.OrganisationID)
				h.mu.Unlock()
				h.broadcastPresence(client.OrganisationID)
			} else {
				h.mu.Unlock()
			}
		}
	}
}

// GetOnlineUsers returns a list of online user IDs in an organisation
func (h *Hub) GetOnlineUsers(orgID uint) []uint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userSet := make(map[uint]bool)
	for client := range h.Clients {
		if client.OrganisationID == orgID {
			userSet[client.UserID] = true
		}
	}

	online := make([]uint, 0, len(userSet))
	for userID := range userSet {
		online = append(online, userID)
	}
	return online
}

// broadcastPresence broadcasts the list of online users to all clients in the organisation
func (h *Hub) broadcastPresence(orgID uint) {
	onlineUsers := h.GetOnlineUsers(orgID)
	msg, err := json.Marshal(PresenceMessage{OnlineUsers: onlineUsers})
	if err != nil {
		log.Printf("Error marshaling presence message: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.Clients {
		if client.OrganisationID == orgID {
			select {
			case client.Send <- msg:
			default:
				// If send buffer is full, skip
			}
		}
	}
}

// BroadcastToUser broadcasts a message to all active sessions of a specific user
func (h *Hub) BroadcastToUser(userID uint, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.Clients {
		if client.UserID == userID {
			select {
			case client.Send <- message:
			default:
				// If send buffer is full, skip
			}
		}
	}
}