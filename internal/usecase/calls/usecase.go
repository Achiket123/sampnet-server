package calls

import (
	"context"
	"log"
	"sync"
	domain "server/internal/domain/calls"
	"github.com/gorilla/websocket"
)

type Room struct {
	Offer      *Message
	Clients    map[*websocket.Conn]bool
	Candidates []Message
	mu         sync.Mutex
}

type Message struct {
	Type      string      `json:"type"`
	Offer     interface{} `json:"offer,omitempty"`
	Answer    interface{} `json:"answer,omitempty"`
	Candidate interface{} `json:"candidate,omitempty"`
}

var rooms = make(map[string]*Room)
var roomsMu sync.Mutex

type service struct{}

func NewService() domain.UseCase {
	return &service{}
}

func (s *service) HandleRoom(ctx context.Context, roomID string, ws *websocket.Conn) {
	log.Printf("Client joined room: %s", roomID)

	roomsMu.Lock()
	room, ok := rooms[roomID]
	if !ok {
		room = &Room{
			Clients:    make(map[*websocket.Conn]bool),
			Candidates: []Message{},
		}
		rooms[roomID] = room
	}
	roomsMu.Unlock()

	if room.Offer != nil {
		log.Print("Sending offer to client")
		ws.WriteJSON(*room.Offer)
	}

	room.mu.Lock()
	room.Clients[ws] = true
	room.mu.Unlock()

	defer func() {
		room.mu.Lock()
		delete(room.Clients, ws)
		if len(room.Clients) == 0 {
			roomsMu.Lock()
			delete(rooms, roomID)
			roomsMu.Unlock()
		}
		room.mu.Unlock()
		ws.Close()
	}()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from room %s: %v", roomID, err)
			break
		}

		switch msg.Type {
		case "offer":
			s.handleOffer(room, ws, msg)
		case "answer":
			s.handleAnswer(room, ws, msg)
		case "candidate":
			s.handleCandidate(room, msg)
		case "endCall":
			s.handleEndCall(room, ws)
		}
	}
}

func (s *service) handleOffer(room *Room, ws *websocket.Conn, msg Message) {
	room.mu.Lock()
	defer room.mu.Unlock()

	room.Offer = &msg
	for client := range room.Clients {
		if client != ws {
			client.WriteJSON(msg)
		}
	}
}

func (s *service) handleAnswer(room *Room, ws *websocket.Conn, msg Message) {
	room.mu.Lock()
	defer room.mu.Unlock()

	for client := range room.Clients {
		if client != ws {
			client.WriteJSON(msg)
		}
	}
}

func (s *service) handleCandidate(room *Room, msg Message) {
	room.mu.Lock()
	defer room.mu.Unlock()

	if len(room.Clients) < 2 {
		room.Candidates = append(room.Candidates, msg)
		return
	}

	for client := range room.Clients {
		client.WriteJSON(msg)
	}
}

func (s *service) handleEndCall(room *Room, ws *websocket.Conn) {
	room.mu.Lock()
	defer room.mu.Unlock()

	for client := range room.Clients {
		if client != ws {
			client.WriteJSON(Message{Type: "endCall"})
		}
	}
}
