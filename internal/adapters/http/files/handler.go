package files

import (
	"io"
	"net/http"
	domain "server/internal/domain/files"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileDomain := &domain.File{
		FileName: c.PostForm("file_name"),
		FileType: c.PostForm("file_type"),
		FileSize: int64(len(imageData)),
		Data:     imageData,
	}

	if err := h.uc.UploadFile(c.Request.Context(), fileDomain); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_id": fileDomain.ID, "url": fileDomain.URL})
}

func (h *Handler) GetFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	file, err := h.uc.GetFile(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.Redirect(http.StatusFound, file.URL)
}
