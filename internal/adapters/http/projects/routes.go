package projects

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	projectGroup := r.Group("/api/v1/projects")
	projectGroup.Use(validateToken)
	{
		projectGroup.POST("/create", h.CreateProject)
		projectGroup.POST("", h.CreateProject)
		projectGroup.GET("/:id", h.GetProject)
		projectGroup.PUT("/:id", h.UpdateProject)
		projectGroup.DELETE("/:id", h.DeleteProject)
		projectGroup.GET("/organisation/:organisation_id", h.GetProjectsByOrganisation)
		projectGroup.GET("", h.GetProjectsByOrganisation)
		projectGroup.GET("/team/:team_id", h.GetProjectsByTeam)
		projectGroup.GET("/less-data/:organisation_id", h.GetProjectsWithLessData)
		projectGroup.GET("/less-data", h.GetProjectsWithLessData)
		
		// Milestones
		projectGroup.POST("/:id/milestones", h.CreateMilestone)
		projectGroup.PUT("/:id/milestones/:milestone_id", h.UpdateMilestone)
		projectGroup.DELETE("/:id/milestones/:milestone_id", h.DeleteMilestone)
	}
}
