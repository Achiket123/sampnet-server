package callstate

import (
	"net/http"
	domain "server/internal/domain/callstate"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ uc domain.UseCase }

func NewHandler(uc domain.UseCase) *Handler { return &Handler{uc: uc} }

func (h *Handler) Upsert(c *gin.Context) {
	var state domain.State
	if err := c.BindJSON(&state); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.uc.CreateOrUpdate(c.Request.Context(), &state); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save call state", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"call": state})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	state, err := h.uc.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Call state not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"call": state})
}

func (h *Handler) Offer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	var req struct {
		CallingID        string `json:"calling_id"`
		CallingFirstName string `json:"calling_first_name"`
		CallingLastName  string `json:"calling_last_name"`
		Offer            string `json:"offer"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if err := h.uc.CreateOffer(c.Request.Context(), uint(id), req.CallingID, req.CallingFirstName, req.CallingLastName, req.Offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create offer", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "offer updated"})
}

func (h *Handler) End(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	if err := h.uc.EndCall(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end call", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "call ended"})
}
