package chats

import (
	"net/http"
	domain "server/internal/domain/chats"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ uc domain.UseCase }

func NewHandler(uc domain.UseCase) *Handler { return &Handler{uc: uc} }

func (h *Handler) Create(c *gin.Context) {
	var chat domain.Chat
	if err := c.BindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.uc.CreateChat(c.Request.Context(), &chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"chat": chat})
}

func (h *Handler) List(c *gin.Context) {
	orgID, err := strconv.Atoi(c.Query("organisation_id"))
	if err != nil || orgID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organisation_id is required"})
		return
	}
	chats, err := h.uc.GetChats(c.Request.Context(), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"chats": chats})
}
