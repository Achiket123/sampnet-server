package notifications

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	group := r.Group("/api/v1/notifications")
	group.Use(validateToken)
	{
		group.GET("/", h.GetNotifications)
		group.PUT("/:notification_id/read", h.MarkRead)
		group.PUT("/read-all", h.MarkAllRead)
	}

	// Register WebSocket upgrade endpoint without validateToken middleware
	// as query parameter authentication is handled internally.
	r.GET("/api/v1/ws", h.Upgrade)
}

