package calendar

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	group := r.Group("/api/v1/calendar")
	group.Use(validateToken)
	{
		group.GET("/events", h.GetEvents)
	}
}
