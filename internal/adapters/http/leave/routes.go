package leave

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	leaveRoutes := r.Group("/api/v1/leave")
	leaveRoutes.Use(validateToken)
	{
		leaveRoutes.POST("/request", h.RequestLeave)
		leaveRoutes.GET("/my", h.GetMyLeaves)
		leaveRoutes.GET("/organisation/:organisation_id", h.GetOrganisationLeaves)
		leaveRoutes.GET("/:id", h.GetLeaveByID)
		leaveRoutes.PUT("/:leave_id/approve", h.ApproveLeave)
		leaveRoutes.PUT("/:leave_id/reject", h.RejectLeave)
		leaveRoutes.PUT("/:leave_id/cancel", h.CancelLeave)
	}
}
