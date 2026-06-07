package attendance

import (
	"net/http"
	domain "server/internal/domain/attendance"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc domain.UseCase
}

func NewHandler(uc domain.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) PostAttendance(c *gin.Context) {
	var att domain.Attendance
	if err := c.ShouldBindJSON(&att); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.RecordAttendance(c.Request.Context(), &att); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attendance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Attendance recorded successfully", "attendance": att})
}

func (h *Handler) UpdateAttendance(c *gin.Context) {
	var req struct {
		UserID        uint       `json:"user_id"`
		CheckOutTime  *time.Time `json:"check_out_time"`
		CheckOutPhoto int        `json:"check_out_photo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	att, err := h.uc.UpdateAttendance(c.Request.Context(), req.UserID, req.CheckOutTime, req.CheckOutPhoto)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully", "attendance": att})
}

func (h *Handler) GetAttendanceByDateAndUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	att, err := h.uc.GetAttendanceByDateAndUser(c.Request.Context(), uint(userID), time.Now())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance fetched successfully", "attendance": att})
}

func (h *Handler) GetEmployeeAttendanceHistory(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

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

	history, err := h.uc.GetEmployeeAttendanceHistory(c.Request.Context(), uint(userID), fromTime, toTime, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
