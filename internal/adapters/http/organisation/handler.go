package organisation

import (
	"errors"
	"log"
	"net/http"
	domain "server/internal/domain/organisation"
	usecase "server/internal/usecase/organisation"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler exposes HTTP handlers for organisation operations.
type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) RegisterOrganisation(c *gin.Context) {
	var org domain.Entity
	if err := c.BindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	ownerUserID := userIDVal.(uint)

	employee, err := h.uc.Register(c.Request.Context(), &org, ownerUserID)
	if err != nil {
		log.Default().Println("Failed to register organisation:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organisation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Organisation created successfully", "organisation": org, "employee": employee})
}

func (h *Handler) GetOrganisation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	org, err := h.uc.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organisation not found"})
			return
		}
		if errors.Is(err, usecase.ErrInvalidID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organisation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organisation": org})
}

func (h *Handler) UpdateOrganisation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	var org domain.Entity
	if err := c.BindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	org.ID = uint(id)

	if err := h.uc.Update(c.Request.Context(), &org); err != nil {
		if errors.Is(err, usecase.ErrInvalidID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organisation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organisation": org})
}
