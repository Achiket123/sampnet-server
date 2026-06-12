package settings

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	settingsRoutes := r.Group("/api/v1/settings")
	settingsRoutes.Use(validateToken)
	{
		settingsRoutes.GET("/organisation", h.GetOrgSettings)
		settingsRoutes.PUT("/organisation", h.UpdateOrgSettings)
		settingsRoutes.DELETE("/organisation", h.DeleteOrganisation)

		settingsRoutes.GET("/plan", h.GetPlanInfo)

		settingsRoutes.GET("/role-permissions", h.GetRolePermissions)
		settingsRoutes.PUT("/role-permissions", h.UpdateRolePermissions)

		settingsRoutes.GET("/attendance-policy", h.GetAttendancePolicy)
		settingsRoutes.PUT("/attendance-policy", h.UpdateAttendancePolicy)

		settingsRoutes.GET("/leave-policy", h.GetLeavePolicies)
		settingsRoutes.PUT("/leave-policy", h.UpdateLeavePolicies)

		settingsRoutes.GET("/task-types", h.GetTaskTypes)
		settingsRoutes.POST("/task-types", h.SaveTaskType)
		settingsRoutes.PUT("/task-types/:id", h.SaveTaskType)
		settingsRoutes.DELETE("/task-types/:id", h.DeleteTaskType)

		settingsRoutes.GET("/profile", h.GetUserProfile)
		settingsRoutes.PUT("/profile", h.UpdateUserProfile)
		settingsRoutes.PUT("/change-password", h.UpdatePassword)

		settingsRoutes.GET("/notification-preferences", h.GetNotificationPreferences)
		settingsRoutes.PUT("/notification-preferences", h.UpdateNotificationPreferences)

		settingsRoutes.GET("/export", h.ExportOrgData)
	}
}
