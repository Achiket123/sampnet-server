package messages

import (
	"net/http"
	domain "server/internal/domain/messages"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ uc domain.UseCase }

func NewHandler(uc domain.UseCase) *Handler { return &Handler{uc: uc} }

func (h *Handler) GetMessages(c *gin.Context) {
	peerID := c.Param("peer_id")
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	msgs, err := h.uc.GetMessages(c.Request.Context(), userIDVal.(uint), peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

func (h *Handler) SendMessage(c *gin.Context) {
	var msg domain.Message
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if _, err := strconv.Atoi(msg.SenderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender_id must be numeric"})
		return
	}
	if _, err := strconv.Atoi(msg.ReceiverID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receiver_id must be numeric"})
		return
	}
	if err := h.uc.SendMessage(c.Request.Context(), &msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": msg})
}
