package messages

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	g := r.Group("/api/v1/messages")
	g.Use(validateToken)
	{
		g.GET("/:peer_id", h.GetMessages)
		g.POST("/send", h.SendMessage)
	}
}
