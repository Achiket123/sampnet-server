package organisation

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	organisationGroup := r.Group("/api/v1/organisation")
	organisationGroup.Use(validateToken)
	{
		organisationGroup.POST("/register", h.RegisterOrganisation)
		organisationGroup.GET("/get/:id", h.GetOrganisation)
		organisationGroup.PUT("/update/:id", h.UpdateOrganisation)
	}
}
