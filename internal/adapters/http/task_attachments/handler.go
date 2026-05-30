package task_attachments

import (
	"net/http"
	"server/internal/domain/task_attachments"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc task_attachments.UseCase
}

func NewHandler(uc task_attachments.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) AddAttachment(c *gin.Context) {
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
		FileID   uint   `json:"file_id" binding:"required"`
		FileName string `json:"file_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID and File Name are required"})
		return
	}

	attachment, err := h.uc.AttachFile(c.Request.Context(), uint(taskID), req.FileID, userIDVal.(uint), req.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attachment"})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func (h *Handler) GetAttachments(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	attachments, err := h.uc.GetAttachments(c.Request.Context(), uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attachments"})
		return
	}

	c.JSON(http.StatusOK, attachments)
}

func (h *Handler) RemoveAttachment(c *gin.Context) {
	attachmentID, err := strconv.Atoi(c.Param("attachment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.uc.RemoveAttachment(c.Request.Context(), uint(attachmentID), userIDVal.(uint))
	if err != nil {
		if err == task_attachments.ErrNotAttachmentOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err == task_attachments.ErrAttachmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove attachment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attachment removed successfully"})
}
