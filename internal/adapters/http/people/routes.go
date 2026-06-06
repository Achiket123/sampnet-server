package people

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	peopleRoutes := r.Group("/api/v1/people")
	peopleRoutes.Use(validateToken)
	{
		// Contacts
		peopleRoutes.GET("/contacts", h.GetContacts)
		peopleRoutes.POST("/contacts", h.CreateContact)
		peopleRoutes.GET("/contacts/:id", h.GetContactByID)
		peopleRoutes.PUT("/contacts/:id", h.UpdateContact)
		peopleRoutes.DELETE("/contacts/:id", h.DeleteContact)
		peopleRoutes.POST("/contacts/bulk-stage", h.BulkUpdateStage)

		// Interactions
		peopleRoutes.POST("/contacts/:id/interactions", h.AddInteraction)

		// Analytics
		peopleRoutes.GET("/analytics", h.GetAnalytics)

		// Lists
		peopleRoutes.GET("/lists", h.GetLists)
		peopleRoutes.POST("/lists", h.CreateList)

		// Pipeline Stages
		peopleRoutes.GET("/stages", h.GetPipelineStages)
		peopleRoutes.POST("/stages", h.CreatePipelineStage)
		peopleRoutes.POST("/stages/reorder", h.ReorderStages)
	}
}
