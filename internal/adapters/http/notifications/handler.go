package notifications

import (
	"net/http"
	"server/internal/domain/notifications"
	"server/internal/platform/websocket"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc  notifications.UseCase
	hub *websocket.Hub
}

func NewHandler(uc notifications.UseCase, hub *websocket.Hub) *Handler {
	return &Handler{
		uc:  uc,
		hub: hub,
	}
}

func (h *Handler) GetNotifications(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	ns, err := h.uc.GetNotifications(c.Request.Context(), userID, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": ns})
}

func (h *Handler) MarkRead(c *gin.Context) {
	idStr := c.Param("notification_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	err = h.uc.MarkNotificationRead(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

func (h *Handler) MarkAllRead(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	err := h.uc.MarkAllNotificationsRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}
