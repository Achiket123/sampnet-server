package callstate

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	g := r.Group("/api/v1/calls")
	g.Use(validateToken)
	{
		g.POST("/upsert", h.Upsert)
		g.GET("/:id", h.Get)
		g.PUT("/offer/:id", h.Offer)
		g.PUT("/end/:id", h.End)
	}
}
