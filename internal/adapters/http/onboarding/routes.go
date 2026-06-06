package onboarding

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	group := r.Group("/api/v1/onboarding")
	group.Use(validateToken)
	{
		group.GET("/:user_id", h.GetOnboardingProgress)
		group.POST("/update", h.UpdateOnboardingProgress)
	}
}
