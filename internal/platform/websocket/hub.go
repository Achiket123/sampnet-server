package websocket

import (
	"log"
	"sync"
)

// PresenceMessage is sent to clients when the list of online users changes
type PresenceMessage struct {
	OnlineUsers []uint `json:"online_users"`
}

type HubMessage struct {
	TargetUserID string
	TargetRoomID string
	SenderClient *Client
	Payload      []byte
}

type RoomSubscription struct {
	Client *Client
	RoomID string
}

type Hub struct {
	// Registered clients mapped by userID (string)
	Clients map[string]map[*Client]bool

	// Subscribed clients mapped by roomID (string)
	Rooms map[string]map[*Client]bool

	Register   chan *Client
	Unregister chan *Client
	Subscribe  chan *RoomSubscription
	Broadcast  chan *HubMessage

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]map[*Client]bool),
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Subscribe:  make(chan *RoomSubscription),
		Broadcast:  make(chan *HubMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.UserID] == nil {
				h.Clients[client.UserID] = make(map[*Client]bool)
			}
			h.Clients[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered: user %s", client.UserID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Clients[client.UserID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					log.Printf("Client unregistered: user %s", client.UserID)
				}
				if len(clients) == 0 {
					delete(h.Clients, client.UserID)
				}
			}
			// Remove from all rooms
			for roomID, roomClients := range h.Rooms {
				if _, exists := roomClients[client]; exists {
					delete(roomClients, client)
					if len(roomClients) == 0 {
						delete(h.Rooms, roomID)
					}
				}
			}
			h.mu.Unlock()

		case sub := <-h.Subscribe:
			h.mu.Lock()
			if h.Rooms[sub.RoomID] == nil {
				h.Rooms[sub.RoomID] = make(map[*Client]bool)
			}
			h.Rooms[sub.RoomID][sub.Client] = true
			h.mu.Unlock()
			log.Printf("Client user %s subscribed to room %s", sub.Client.UserID, sub.RoomID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			var targetClients []*Client

			if message.TargetUserID != "" {
				if clients, ok := h.Clients[message.TargetUserID]; ok {
					for c := range clients {
						targetClients = append(targetClients, c)
					}
				}
			}

			if message.TargetRoomID != "" {
				if clients, ok := h.Rooms[message.TargetRoomID]; ok {
					for c := range clients {
						targetClients = append(targetClients, c)
					}
				}
			}

			// Deduplicate clients and exclude sender
			uniqueClients := make(map[*Client]bool)
			for _, c := range targetClients {
				if message.SenderClient != nil && c == message.SenderClient {
					continue
				}
				uniqueClients[c] = true
			}

			// Send to collected unique clients
			for client := range uniqueClients {
				select {
				case client.Send <- message.Payload:
				default:
					h.mu.RUnlock()
					h.Unregister <- client
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// IsOnline checks if a user is currently connected
func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.Clients[userID]
	return ok
}