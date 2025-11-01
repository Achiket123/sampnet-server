package calls
import (
	"log"
	"net/http"
	"sync" 
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Room represents a room for signaling
type Room struct {
	Offer       *Message          // Stores the offer from the first participant
	Clients     map[*websocket.Conn]bool // Connected clients
	Candidates  []Message          // Stores ICE candidates for delayed processing
	mu          sync.Mutex         // Ensures thread-safe operations
}

// Global map to manage rooms
var rooms = make(map[string]*Room)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow requests from any origin because token handling is done through middleware
		return true
	},
}

// Message defines the structure of WebRTC signaling messages
type Message struct {
	Type      string      `json:"type"`                 // "offer", "answer", "candidate"
	Offer     interface{} `json:"offer,omitempty"`     // SDP offer
	Answer    interface{} `json:"answer,omitempty"`    // SDP answer
	Candidate interface{} `json:"candidate,omitempty"` // ICE candidate
}

// Handle WebSocket connections for a specific room
func HandleRoom(c *gin.Context) {
	roomID := c.Param("id")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer ws.Close()
 
	log.Printf("Client joined room: %s", roomID)

	// Get or create the room
	room, ok := rooms[roomID]
	if !ok {
		// Create a new room if it doesn't exist
		room = &Room{
			Clients:    make(map[*websocket.Conn]bool),
			Candidates: []Message{},
		}
		rooms[roomID] = room
	}
	if room.Offer != nil {
		// Send the offer to the client
		log.Print("Sending offer to client") 
		err := ws.WriteJSON(*room.Offer)
		if err != nil {
			log.Println("Error sending offer:", err)
			return
		}
	}


	room.mu.Lock()
	room.Clients[ws] = true
	room.mu.Unlock()

	// Listen for messages from this client
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from room %s: %v", roomID, err)
			room.mu.Lock()
			delete(room.Clients, ws)
			delete(rooms, roomID)
			room.mu.Unlock()
			break
		}

		switch msg.Type {
		case "offer":
			HandleOffer(room, ws, msg)
		case "answer":
			 
			HandleAnswer(room, ws, msg)
		case "candidate":
			HandleCandidate(room, msg)
		case "endCall":
			room.mu.Lock()
			// Notify all other clients in the room
			for client := range room.Clients {
				if client != ws {
					err := client.WriteJSON(Message{Type: "endCall"})
					if err != nil {
						log.Println("Error sending endCall message:", err)
						client.Close()
						delete(room.Clients, client)
					}
				}
			}
			room.mu.Unlock()
		
			// (Optional) If no clients remain, delete the room
			if len(room.Clients) == 0 {
				delete(rooms, roomID)
				log.Printf("Room %s deleted due to no clients.", roomID)
			}
		
		} 
	}
}

// Handle incoming offer from the first participant
func HandleOffer(room *Room, ws *websocket.Conn, msg Message) {
	room.mu.Lock()
	defer room.mu.Unlock()

	// Store the offer in the room
	room.Offer = &msg
	log.Println("Offer received and stored for room.")

	// Notify all clients except the sender about the offer
	for client := range room.Clients {
		if client != ws {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error sending offer to client:", err)
				client.Close()
				delete(room.Clients, client)
			}
		}
	}
}

// Handle incoming answer from the second participant
func HandleAnswer(room *Room, ws *websocket.Conn, msg Message) {
	room.mu.Lock()
	defer room.mu.Unlock()

	// Notify all clients except the sender about the answer
	for client := range room.Clients {
		if client != ws {
			log.Print("Sending answer to client")
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error sending answer to client:", err)
				client.Close()
				delete(room.Clients, client)
			}
		}
	}
}

// Handle incoming ICE candidates
func HandleCandidate(room *Room, msg Message) {
	room.mu.Lock()
	defer room.mu.Unlock()

	// Store ICE candidates if no remote peer exists yet
	if len(room.Clients) < 2 {
		room.Candidates = append(room.Candidates, msg)
		return
	}

	// Broadcast ICE candidates to all connected clients
	for client := range room.Clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Println("Error sending ICE candidate to client:", err)
			client.Close()
			delete(room.Clients, client)
		}
	}
}

func HandleDisconnect(roomID string) {
	room, ok := rooms[roomID]
	if !ok {
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	for client := range room.Clients {	
		err := client.WriteJSON(Message{Type: "endCall"})
		if err != nil {	
			log.Println("Error sending endCall message:", err)
			client.Close()
			delete(room.Clients, client)
		}	
	}
	if( len(room.Clients) == 0 ){
		delete(rooms, roomID)
	}
	}