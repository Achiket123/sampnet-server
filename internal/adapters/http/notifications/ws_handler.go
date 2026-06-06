package notifications

import (
	"fmt"
	"log"
	"net/http"
	"server/internal/platform/miscallenous"
	"server/internal/platform/websocket"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) Upgrade(c *gin.Context) {
	// Extract the token from the "token" query parameter
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Decode and validate the token
	decodedToken, err := miscallenous.DecodeJWTToken(tokenString)
	if err != nil || !decodedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			userID = uint(sub)
		}
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	// Upgrade the connection to WebSocket
	conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &websocket.Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: fmt.Sprintf("%d", userID),
	}

	client.Hub.Register <- client

	// Start read/write pumps in goroutines
	go client.WritePump()
	go client.ReadPump()
}
