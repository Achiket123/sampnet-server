package task_comments

import (
	"net/http"
	"server/internal/domain/task_comments"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc task_comments.UseCase
}

func NewHandler(uc task_comments.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) AddComment(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	comment, err := h.uc.AddComment(c.Request.Context(), uint(taskID), userIDVal.(uint), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) GetComments(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	comments, err := h.uc.GetComments(c.Request.Context(), uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *Handler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.uc.DeleteComment(c.Request.Context(), uint(commentID), userIDVal.(uint))
	if err != nil {
		if err == task_comments.ErrNotCommentOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err == task_comments.ErrCommentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
