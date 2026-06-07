package leave

import (
	"net/http"
	"server/internal/domain/leave"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc leave.UseCase
}

func NewHandler(uc leave.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) RequestLeave(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, _ := userIDStr.(uint)

	var req struct {
		OrganisationID uint   `json:"organisation_id"`
		LeaveType      string `json:"leave_type"`
		StartDate      string `json:"start_date"`
		EndDate        string `json:"end_date"`
		Reason         string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate leave type
	validTypes := map[string]bool{
		"annual":    true,
		"sick":      true,
		"unpaid":    true,
		"maternity": true,
		"paternity": true,
		"emergency": true,
	}
	if !validTypes[req.LeaveType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid leave type"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use ISO 8601"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use ISO 8601"})
		return
	}

	l, err := h.uc.RequestLeave(c.Request.Context(), userID, req.OrganisationID, req.LeaveType, startDate, endDate, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, l)
}

func (h *Handler) GetMyLeaves(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, _ := userIDStr.(uint)

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	leaves, err := h.uc.GetMyLeaves(c.Request.Context(), userID, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaves": leaves})
}

func (h *Handler) GetOrganisationLeaves(c *gin.Context) {
	orgIDStr := c.Param("organisation_id")
	orgID, _ := strconv.ParseUint(orgIDStr, 10, 32)
	status := c.Query("status")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	leaves, err := h.uc.GetOrganisationLeaves(c.Request.Context(), uint(orgID), status, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaves": leaves})
}

func (h *Handler) GetLeaveByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	l, err := h.uc.GetLeave(c.Request.Context(), uint(id))
	if err != nil {
		if err == leave.ErrLeaveNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leave": l})
}

func (h *Handler) ApproveLeave(c *gin.Context) {
	idStr := c.Param("leave_id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	managerIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	managerID, _ := managerIDStr.(uint)

	var req struct {
		ManagerNote string `json:"manager_note"`
	}
	c.ShouldBindJSON(&req)

	err := h.uc.ApproveLeave(c.Request.Context(), uint(id), managerID, req.ManagerNote)
	if err != nil {
		if err == leave.ErrLeaveNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "leave is not in a pending state" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "leave approved"})
}

func (h *Handler) RejectLeave(c *gin.Context) {
	idStr := c.Param("leave_id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	managerIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	managerID, _ := managerIDStr.(uint)

	var req struct {
		ManagerNote string `json:"manager_note"`
	}
	c.ShouldBindJSON(&req)

	err := h.uc.RejectLeave(c.Request.Context(), uint(id), managerID, req.ManagerNote)
	if err != nil {
		if err == leave.ErrLeaveNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "leave is not in a pending state" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "leave rejected"})
}

func (h *Handler) CancelLeave(c *gin.Context) {
	idStr := c.Param("leave_id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, _ := userIDStr.(uint)

	err := h.uc.CancelLeave(c.Request.Context(), uint(id), userID)
	if err != nil {
		if err == leave.ErrLeaveNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == leave.ErrNotAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "leave cancelled"})
}

func (h *Handler) GetEmployeeLeaveHistory(c *gin.Context) {
	employeeIDStr := c.Param("employeeId")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	status := c.Query("status")
	
	limitStr := c.Query("limit")
	limit := 10 // default limit
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offsetStr := c.Query("offset")
	offset := 0
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	var fromTime *time.Time
	if from := c.Query("from"); from != "" {
		if parsed, err := time.Parse(time.RFC3339, from); err == nil {
			fromTime = &parsed
		} else if parsed, err := time.Parse("2006-01-02", from); err == nil {
			fromTime = &parsed
		}
	}

	var toTime *time.Time
	if to := c.Query("to"); to != "" {
		if parsed, err := time.Parse(time.RFC3339, to); err == nil {
			toTime = &parsed
		} else if parsed, err := time.Parse("2006-01-02", to); err == nil {
			toTime = &parsed
		}
	}

	history, err := h.uc.GetEmployeeLeaveHistory(c.Request.Context(), uint(employeeID), status, fromTime, toTime, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leave history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
