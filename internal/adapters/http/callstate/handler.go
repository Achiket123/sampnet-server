package callstate

import (
	"fmt"
	"net/http"
	domain "server/internal/domain/callstate"
	chatDomain "server/internal/domain/chats"
	ws "server/internal/platform/websocket"
	"strconv"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc        domain.UseCase
	wsManager *ws.Manager
	chatRepo  chatDomain.Repository
}

func NewHandler(uc domain.UseCase, wsManager *ws.Manager, chatRepo chatDomain.Repository) *Handler {
	return &Handler{uc: uc, wsManager: wsManager, chatRepo: chatRepo}
}

func (h *Handler) Upsert(c *gin.Context) {
	var state domain.State
	if err := c.BindJSON(&state); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.uc.CreateOrUpdate(c.Request.Context(), &state); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save call state", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"call": state})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	state, err := h.uc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Call state not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"call": state})
}

func (h *Handler) Offer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	var req struct {
		CallingID        string `json:"calling_id"`
		CallingFirstName string `json:"calling_first_name"`
		CallingLastName  string `json:"calling_last_name"`
		Offer            string `json:"offer"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.uc.CreateOffer(c.Request.Context(), uint(id), req.CallingID, req.CallingFirstName, req.CallingLastName, req.Offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create offer", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer updated"})
}

func (h *Handler) End(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	if err := h.uc.EndCall(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end call", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "call ended"})
}

// InitiateCall sends a call invitation via WebSocket
func (h *Handler) InitiateCall(c *gin.Context) {
	var req struct {
		TargetUserID string `json:"target_user_id"`
		Type         string `json:"type"` // "audio" or "video"
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDVal.(uint)
	userIDStr := fmt.Sprintf("%d", userID)

	orgIDVal, _ := c.Get("organisationID")
	orgID := orgIDVal.(uint)

	targetUserIDUint, err := strconv.ParseUint(req.TargetUserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user id"})
		return
	}

	chat, err := h.chatRepo.GetOrCreateDM(c.Request.Context(), userID, uint(targetUserIDUint), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create DM"})
		return
	}

	callID := uuid.New().String()

	if h.wsManager != nil {
		payload := map[string]interface{}{
			"caller_id": userIDStr,
			"type":      req.Type,
			"room_id":   chat.RoomID,
			"call_id":   callID,
		}
		_ = h.wsManager.SendToUser(req.TargetUserID, "call_incoming", payload)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "call initiated",
		"room_id": chat.RoomID,
		"call_id": callID,
	})
}
