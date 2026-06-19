package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer (increased for WebRTC SDP payloads).
	maxMessageSize = 16384
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local dev
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	// Client metadata.
	UserID string
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Try to parse basic envelope to handle internal commands like subscribe_room
		// and route WebRTC signaling messages directly
		var envelope struct {
			Type          string   `json:"type"`
			RoomID        string   `json:"room_id"` // Payload structure depending on type
			TargetUserIDs []string `json:"target_user_ids"`
		}
		if err := json.Unmarshal(message, &envelope); err == nil {
			if envelope.Type == "subscribe_room" && envelope.RoomID != "" {
				c.Hub.Subscribe <- &RoomSubscription{
					Client: c,
					RoomID: envelope.RoomID,
				}
			} else if envelope.Type == "call_offer" || envelope.Type == "call_answer" || envelope.Type == "ice_candidate" || envelope.Type == "call_ended" || envelope.Type == "call_rejected" || envelope.Type == "call_accepted" {
				hasTargets := false
				// Broadcast to all specified target users
				for _, targetID := range envelope.TargetUserIDs {
					if targetID != "" {
						hasTargets = true
						c.Hub.Broadcast <- &HubMessage{
							TargetUserID: targetID,
							SenderClient: c,
							Payload:      message,
						}
					}
				}
				// Also broadcast to the room only if no specific targets were addressed
				if !hasTargets && envelope.RoomID != "" {
					c.Hub.Broadcast <- &HubMessage{
						TargetRoomID: envelope.RoomID,
						SenderClient: c,
						Payload:      message,
					}
				}
			}
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}