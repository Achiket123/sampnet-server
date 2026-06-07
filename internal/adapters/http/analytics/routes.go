package analytics

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc, roleMiddleware gin.HandlerFunc) {
	analyticsGroup := r.Group("/api/v1/analytics")
	analyticsGroup.Use(validateToken)
	// Apply Boss/Manager role check
	analyticsGroup.Use(roleMiddleware)
	{
		analyticsGroup.GET("/employee/:userId", h.GetEmployeeAnalytics)
		analyticsGroup.GET("/organisation/:orgId/employees", h.GetOrgEmployeeMonitor)
	}
}