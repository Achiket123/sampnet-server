package files

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	fileGroup := r.Group("/api/v1/file")
	{
		fileGroup.POST("/upload", h.UploadFile)
		fileGroup.GET("/:id", h.GetFile)
	}
}
