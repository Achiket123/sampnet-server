package onboarding

import (
	"net/http"
	"server/internal/domain/onboarding"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc onboarding.UseCase
}

func NewHandler(uc onboarding.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetOnboardingProgress(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	p, err := h.uc.GetOnboardingProgress(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *Handler) UpdateOnboardingProgress(c *gin.Context) {
	var p onboarding.OnboardingProgress
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.UpdateOnboardingProgress(c.Request.Context(), &p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "onboarding progress updated"})
}
