package teams

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	teamGroup := r.Group("/api/v1/teams")
	teamGroup.Use(validateToken)
	{
		teamGroup.POST("/create", h.CreateTeam)
		teamGroup.GET("/:id", h.GetTeam)
		teamGroup.PUT("/:id", h.UpdateTeam)
		teamGroup.GET("/organisation/:organisation_id", h.GetTeamsByOrganisation)
		teamGroup.DELETE("/delete/:id", h.DeleteTeam)
	}
}
