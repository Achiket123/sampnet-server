package task_attachments

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	attachments := r.Group("/api/v1/tasks/:task_id/attachments")
	attachments.Use(validateToken)
	{
		attachments.POST("/", h.AddAttachment)
		attachments.GET("/", h.GetAttachments)
		attachments.DELETE("/:attachment_id", h.RemoveAttachment)
	}
}
