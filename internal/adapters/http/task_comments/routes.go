package task_comments

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	comments := r.Group("/api/v1/tasks/:task_id/comments")
	comments.Use(validateToken)
	{
		comments.POST("/", h.AddComment)
		comments.GET("/", h.GetComments)
		comments.DELETE("/:comment_id", h.DeleteComment)
	}
}
