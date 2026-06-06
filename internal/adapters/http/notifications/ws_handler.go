package notifications

import (
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
	var organisationID uint
	var employeeID uint

	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			userID = uint(sub)
		}
		if emp, ok := claims["employee"].(map[string]interface{}); ok {
			if orgID, ok := emp["organisation_id"].(float64); ok {
				organisationID = uint(orgID)
			}
			if id, ok := emp["id"].(float64); ok {
				employeeID = uint(id)
			}
		} else if mgr, ok := claims["manager"].(map[string]interface{}); ok {
			if orgID, ok := mgr["organisation_id"].(float64); ok {
				organisationID = uint(orgID)
			}
			if id, ok := mgr["id"].(float64); ok {
				employeeID = uint(id)
			}
		} else if boss, ok := claims["boss"].(map[string]interface{}); ok {
			if orgID, ok := boss["organisation_id"].(float64); ok {
				organisationID = uint(orgID)
			}
			if id, ok := boss["id"].(float64); ok {
				employeeID = uint(id)
			}
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
		Hub:            h.hub,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		UserID:         userID,
		EmployeeID:     employeeID,
		OrganisationID: organisationID,
	}

	client.Hub.Register <- client

	// Start read/write pumps in goroutines
	go client.WritePump()
	go client.ReadPump()
}
