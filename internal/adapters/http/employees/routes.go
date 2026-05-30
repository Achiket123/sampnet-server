package employees

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *Handler, validateToken gin.HandlerFunc) {
	employeeGroup := r.Group("/api/v1/employees")
	employeeGroup.Use(validateToken)
	{
		employeeGroup.GET("/is-employee/:user_id", h.IsEmployee)
		employeeGroup.POST("/add", h.AddEmployee)
		employeeGroup.GET("/get/:organisation_id", h.GetEmployees)
		employeeGroup.GET("/list/:id", h.GetEmployeeByID)
		employeeGroup.PUT("/update/:id", h.UpdateEmployee)
		employeeGroup.DELETE("/delete/:id", h.DeleteEmployee)
		employeeGroup.GET("/search", h.SearchEmployee)
		employeeGroup.POST("/make-manager", h.MakeManager)
	}

	{
		r.POST("/api/v1/create/boss", validateToken, h.CreateBoss)
		r.GET("/api/v1/get/boss", validateToken, h.GetBoss)
		r.POST("/api/v1/create/employee", validateToken, h.CreateEmployee)
	}
}
