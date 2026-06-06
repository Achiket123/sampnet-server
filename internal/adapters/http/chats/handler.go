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

	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	orgIDVal, _ := c.Get("organisationID")

	chat.CreatedBy = userIDVal.(uint)
	chat.OrganisationID = orgIDVal.(uint)

	if err := h.uc.CreateChat(c.Request.Context(), &chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, chat)
}

func (h *Handler) List(c *gin.Context) {
	orgIDVal, ok := c.Get("organisationID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation not found"})
		return
	}

	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	chats, err := h.uc.GetUserChats(c.Request.Context(), userIDVal.(uint), orgIDVal.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	// if it returns nil, it is empty slice
	if chats == nil {
		chats = []domain.Chat{}
	}
	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

func (h *Handler) GetOrCreateDM(c *gin.Context) {
	peerIDStr := c.Param("peer_id")
	peerID, err := strconv.ParseUint(peerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid peer ID"})
		return
	}

	orgIDVal, ok := c.Get("organisationID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organisation not found"})
		return
	}

	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	chat, err := h.uc.GetOrCreateDM(c.Request.Context(), userIDVal.(uint), uint(peerID), orgIDVal.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get or create DM"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat})
}
