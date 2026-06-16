package calls

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	domain "server/internal/domain/calls"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Room struct {
	Offer      *Message
	Clients    map[*websocket.Conn]bool
	// Buffered ICE candidates are kept in memory on whichever instance first buffered them, and are only replayed to a client that later joins that same instance, cross instance replay of buffered candidates is out of scope for this change.
	Candidates []Message
	mu         sync.Mutex
	PubSub     *redis.PubSub
	CancelSub  context.CancelFunc
}

type Message struct {
	Type      string      `json:"type"`
	Offer     interface{} `json:"offer,omitempty"`
	Answer    interface{} `json:"answer,omitempty"`
	Candidate interface{} `json:"candidate,omitempty"`
}

type callSignalMessage struct {
	OriginInstanceID string      `json:"origin_instance_id"`
	Type             string      `json:"type"`
	Offer            interface{} `json:"offer,omitempty"`
	Answer           interface{} `json:"answer,omitempty"`
	Candidate        interface{} `json:"candidate,omitempty"`
}

var rooms = make(map[string]*Room)
var roomsMu sync.Mutex

type service struct {
	redisClient *redis.Client
	instanceID  string
}

func NewService(redisClient *redis.Client) domain.UseCase {
	return &service{
		redisClient: redisClient,
		instanceID:  uuid.New().String(),
	}
}

func (s *service) HandleRoom(ctx context.Context, roomID string, ws *websocket.Conn) {
	log.Printf("Client joined room: %s", roomID)

	roomsMu.Lock()
	room, ok := rooms[roomID]
	var newlyCreated bool
	if !ok {
		room = &Room{
			Clients:    make(map[*websocket.Conn]bool),
			Candidates: []Message{},
		}
		rooms[roomID] = room
		newlyCreated = true
	}
	roomsMu.Unlock()

	if newlyCreated {
		subCtx, cancelSub := context.WithCancel(context.Background())
		room.CancelSub = cancelSub

		channelName := "call_room:" + roomID
		pubsub := s.redisClient.Subscribe(context.Background(), channelName)
		room.PubSub = pubsub

		go func() {
			ch := pubsub.Channel()
			for {
				select {
				case <-subCtx.Done():
					pubsub.Close()
					return
				case msg, ok := <-ch:
					if !ok {
						return
					}
					var signal callSignalMessage
					if err := json.Unmarshal([]byte(msg.Payload), &signal); err != nil {
						log.Printf("Error unmarshaling call signal message: %v", err)
						continue
					}
					if signal.OriginInstanceID == s.instanceID {
						continue
					}

					room.mu.Lock()
					localMsg := Message{
						Type:      signal.Type,
						Offer:     signal.Offer,
						Answer:    signal.Answer,
						Candidate: signal.Candidate,
					}
					for client := range room.Clients {
						client.WriteJSON(localMsg)
					}
					room.mu.Unlock()
				}
			}
		}()
	}

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
			if room.CancelSub != nil {
				room.CancelSub()
			}
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
			s.handleOffer(roomID, room, ws, msg)
		case "answer":
			s.handleAnswer(roomID, room, ws, msg)
		case "candidate":
			s.handleCandidate(roomID, room, msg)
		case "endCall":
			s.handleEndCall(roomID, room, ws)
		}
	}
}

func (s *service) handleOffer(roomID string, room *Room, ws *websocket.Conn, msg Message) {
	room.mu.Lock()
	room.Offer = &msg
	for client := range room.Clients {
		if client != ws {
			client.WriteJSON(msg)
		}
	}
	room.mu.Unlock()

	signal := callSignalMessage{
		OriginInstanceID: s.instanceID,
		Type:             msg.Type,
		Offer:            msg.Offer,
	}
	payload, err := json.Marshal(signal)
	if err != nil {
		log.Printf("Error marshaling call offer signal: %v", err)
		return
	}
	if err := s.redisClient.Publish(context.Background(), "call_room:"+roomID, payload).Err(); err != nil {
		log.Printf("Error publishing call offer signal: %v", err)
	}
}

func (s *service) handleAnswer(roomID string, room *Room, ws *websocket.Conn, msg Message) {
	room.mu.Lock()
	for client := range room.Clients {
		if client != ws {
			client.WriteJSON(msg)
		}
	}
	room.mu.Unlock()

	signal := callSignalMessage{
		OriginInstanceID: s.instanceID,
		Type:             msg.Type,
		Answer:           msg.Answer,
	}
	payload, err := json.Marshal(signal)
	if err != nil {
		log.Printf("Error marshaling call answer signal: %v", err)
		return
	}
	if err := s.redisClient.Publish(context.Background(), "call_room:"+roomID, payload).Err(); err != nil {
		log.Printf("Error publishing call answer signal: %v", err)
	}
}

func (s *service) handleCandidate(roomID string, room *Room, msg Message) {
	room.mu.Lock()
	if len(room.Clients) < 2 {
		room.Candidates = append(room.Candidates, msg)
		room.mu.Unlock()

		s.publishCandidate(roomID, msg)
		return
	}

	for client := range room.Clients {
		client.WriteJSON(msg)
	}
	room.mu.Unlock()

	s.publishCandidate(roomID, msg)
}

func (s *service) publishCandidate(roomID string, msg Message) {
	signal := callSignalMessage{
		OriginInstanceID: s.instanceID,
		Type:             msg.Type,
		Candidate:        msg.Candidate,
	}
	payload, err := json.Marshal(signal)
	if err != nil {
		log.Printf("Error marshaling call candidate signal: %v", err)
		return
	}
	if err := s.redisClient.Publish(context.Background(), "call_room:"+roomID, payload).Err(); err != nil {
		log.Printf("Error publishing call candidate signal: %v", err)
	}
}

func (s *service) handleEndCall(roomID string, room *Room, ws *websocket.Conn) {
	room.mu.Lock()
	for client := range room.Clients {
		if client != ws {
			client.WriteJSON(Message{Type: "endCall"})
		}
	}
	room.mu.Unlock()

	signal := callSignalMessage{
		OriginInstanceID: s.instanceID,
		Type:             "endCall",
	}
	payload, err := json.Marshal(signal)
	if err != nil {
		log.Printf("Error marshaling call end signal: %v", err)
		return
	}
	if err := s.redisClient.Publish(context.Background(), "call_room:"+roomID, payload).Err(); err != nil {
		log.Printf("Error publishing call end signal: %v", err)
	}
}
