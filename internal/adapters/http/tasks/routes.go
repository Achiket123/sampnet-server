package tasks

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	taskGroup := r.Group("/api/v1/tasks")
	taskGroup.Use(validateToken)
	{
		taskGroup.POST("/create", h.CreateTask)
		taskGroup.PUT("/update/:id", h.UpdateTask)
		taskGroup.DELETE("/delete/:id", h.SoftDeleteTask)
		taskGroup.GET("/get/:id", h.GetTaskByID)
		taskGroup.GET("/team", h.GetTeamTasks)
		taskGroup.GET("/project", h.GetProjectTasks)
		taskGroup.GET("/personal", h.GetPersonalTasks)
		taskGroup.GET("/organisation/:organisation_id", h.GetOrganisationTasks)
		taskGroup.GET("/title/:title", h.GetTaskByTitle)
	}
}
