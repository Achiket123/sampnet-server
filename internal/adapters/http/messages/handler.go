package messages

import (
	"fmt"
	"net/http"
	domain "server/internal/domain/messages"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ uc domain.UseCase }

func NewHandler(uc domain.UseCase) *Handler { return &Handler{uc: uc} }

func (h *Handler) GetMessages(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}

	_, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	cursor := c.Query("cursor")
	limitStr := c.Query("limit")
	limit := 30
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	page, err := h.uc.GetMessages(c.Request.Context(), roomID, cursor, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	
	c.JSON(http.StatusOK, page)
}

func (h *Handler) SendMessage(c *gin.Context) {
	var msg domain.Message
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	senderIDStr := fmt.Sprintf("%d", userIDVal.(uint))

	orgIDVal, _ := c.Get("organisationID")
	orgID := orgIDVal.(uint)

	// Enforce sender id and organisation
	msg.SenderID = senderIDStr
	msg.OrganisationID = orgID

	createdMsg, err := h.uc.SendMessage(c.Request.Context(), &msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdMsg)
}

func (h *Handler) MarkSeen(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}

	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDStr := fmt.Sprintf("%d", userIDVal.(uint))

	if err := h.uc.MarkSeen(c.Request.Context(), roomID, userIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark seen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Marked seen"})
}

func (h *Handler) DeleteMessage(c *gin.Context) {
	msgIDStr := c.Param("id")
	msgID, err := strconv.ParseUint(msgIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message id"})
		return
	}

	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userIDStr := fmt.Sprintf("%d", userIDVal.(uint))

	if err := h.uc.DeleteMessage(c.Request.Context(), uint(msgID), userIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}
