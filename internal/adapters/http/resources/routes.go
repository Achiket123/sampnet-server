package resources

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	routes := r.Group("/api/v1/resources")
	routes.Use(validateToken)
	{
		// Collections
		routes.POST("/collections", h.CreateCollection)
		routes.GET("/collections/:id", h.GetCollection)
		routes.GET("/collections", h.ListCollections)
		routes.PUT("/collections/:id", h.UpdateCollection)
		routes.DELETE("/collections/:id", h.DeleteCollection)
		routes.POST("/collections/:id/fields", h.AddField)
		routes.PUT("/collections/:id/fields/:key", h.UpdateField)
		routes.DELETE("/collections/:id/fields/:key", h.RemoveField)

		// Records
		routes.POST("/collections/:id/records", h.CreateRecord)
		routes.GET("/collections/:id/records/:record_id", h.GetRecord)
		routes.GET("/collections/:id/records", h.ListRecords)
		routes.PUT("/collections/:id/records/:record_id", h.UpdateRecord)
		routes.DELETE("/collections/:id/records/:record_id", h.DeleteRecord)
		routes.POST("/collections/:id/bulk", h.BulkCreate)
		routes.GET("/collections/:id/export", h.Export)
	}
}
