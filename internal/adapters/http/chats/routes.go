package chats

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	g := r.Group("/api/v1/chats")
	g.Use(validateToken)
	{
		g.POST("/create", h.Create)
		g.GET("", h.List)
		g.GET("/dm/:peer_id", h.GetOrCreateDM)
	}
}
