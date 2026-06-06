package messages

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	g := r.Group("/api/v1/messages")
	g.Use(validateToken)
	{
		g.GET("/:room_id", h.GetMessages)
		g.POST("/send", h.SendMessage)
		g.PUT("/:room_id/seen", h.MarkSeen)
		g.DELETE("/:id", h.DeleteMessage)
	}
}
