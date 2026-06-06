package search

import "github.com/gin-gonic/gin"

// RegisterRoutes wires the search handler to the Gin engine.
func RegisterRoutes(r *gin.Engine, h *Handler, authMiddleware gin.HandlerFunc) {
	r.GET("/api/v1/search", authMiddleware, h.Search)
}
