package research

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	research := r.Group("/api/v1/research")
	research.Use(validateToken)
	{
		research.POST("", h.Create)
		research.GET("", h.List)
		research.GET("/:id", h.GetByID)
		research.PUT("/:id", h.Update)
		research.DELETE("/:id", h.Delete)

		// Folder Routes
		research.POST("/folders", h.CreateFolder)
		research.GET("/folders", h.ListFolders) // requires ?research_id=X
		research.GET("/folders/:id", h.GetFolder)
		research.PUT("/folders/:id", h.UpdateFolder)
		research.DELETE("/folders/:id", h.DeleteFolder)

		// Document Routes
		research.POST("/documents", h.CreateDocument)
		research.GET("/documents", h.ListDocuments) // requires ?research_id=X
		research.GET("/documents/:id", h.GetDocument)
		research.PUT("/documents/:id", h.UpdateDocument)
		research.DELETE("/documents/:id", h.DeleteDocument)

		// Artifact Routes (Scoped to Documents or Research)
		research.POST("/files", h.UploadFile)
		research.GET("/files", h.ListFiles) // requires ?document_id=X
		research.GET("/files/:id/download", h.DownloadFile)
		research.DELETE("/files/:id", h.DeleteFile)

		research.POST("/references", h.AddReference)
		research.GET("/references", h.ListReferences) // requires ?document_id=X
		research.DELETE("/references/:id", h.DeleteReference)
	}
}
