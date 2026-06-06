package calendar

import (
	"net/http"
	domain "server/internal/domain/calendar"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) GetEvents(c *gin.Context) {
	userIDVal, existsUser := c.Get("userID")
	orgIDVal, existsOrg := c.Get("organisationID")
	if !existsUser || !existsOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "User or organisation not found in context"})
		return
	}

	userID, okUser := userIDVal.(uint)
	orgID, okOrg := orgIDVal.(uint)
	if !okUser || !okOrg {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Invalid user or organisation ID in context"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	viewType := c.DefaultQuery("view_type", "personal")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "start_date and end_date query parameters are required"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		startDate, err = time.Parse("2006-01-02T15:04:05", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid start_date format, must be RFC3339"})
			return
		}
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		endDate, err = time.Parse("2006-01-02T15:04:05", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "invalid end_date format, must be RFC3339"})
			return
		}
	}

	var reqUserID *uint
	if viewType == "personal" {
		reqUserID = &userID
	}

	req := domain.GetEventsRequest{
		OrgID:     orgID,
		UserID:    reqUserID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	events, err := h.uc.GetEvents(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}