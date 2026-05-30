package calls

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/api/v1/call/:id", h.HandleRoom)
}
