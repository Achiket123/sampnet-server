package attendance

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	attendanceGroup := r.Group("/api/v1/attendence")
	attendanceGroup.Use(validateToken)
	{
		attendanceGroup.POST("/create", h.PostAttendance)
		attendanceGroup.PUT("/:id", h.UpdateAttendance)
		attendanceGroup.GET("/:id", h.GetAttendanceByDateAndUser)
	}
}
