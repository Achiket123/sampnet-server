package analytics

import (
	"net/http"
	"server/internal/domain/analytics"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc analytics.UseCase
}

func NewHandler(uc analytics.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetEmployeeAnalytics(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	orgIDStr := c.Query("organisation_id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	period := c.Query("period")
	if period == "" {
		period = "month" // Default
	}

	summary, err := h.uc.GetEmployeeAnalytics(c.Request.Context(), uint(userID), uint(orgID), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *Handler) GetOrgEmployeeMonitor(c *gin.Context) {
	orgIDStr := c.Param("orgId")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation ID"})
		return
	}

	monitorData, err := h.uc.GetOrgEmployeeMonitor(c.Request.Context(), uint(orgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch monitor data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, monitorData)
}
