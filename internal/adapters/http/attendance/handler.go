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
