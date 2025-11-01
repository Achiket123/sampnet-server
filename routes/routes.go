package routes

import (
	"server/calls"
	"server/features/attendence"
	"server/features/auth"
	"server/features/employees"
	"server/features/files"
	"server/features/organisations"
	"server/features/projects"
	"server/features/tasks"
	"server/features/teams"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the API routes for the application
func SetupRoutes(r *gin.Engine) {
	// Authentication routes
	// User authentication routes
	r.POST("/api/v1/auth/signup", auth.SignUp) // Handle user registration
	r.POST("/api/v1/auth/signin", auth.SignIn) // Handle user login
	authGroup := r.Group("/api/v1/auth")
	authGroup.Use(middlewares.ValidateToken())

	{
		authGroup.POST("/complete-signin", auth.CompleteSignIn)    // Handle completing user login
		authGroup.GET("/validate-employee", auth.ValidateEmployee) // Handle validating employee
	}
	// File handling routes
	fileGroup := r.Group("/api/v1/file")
	{
		fileGroup.POST("/upload", files.UploadFile) // Handle file uploads
		fileGroup.GET("/:id", files.GetFile)        // Retrieve a specific file by ID
	}

	// Organization management routes
	organisationGroup := r.Group("/api/v1/organisation")
	organisationGroup.Use(middlewares.ValidateToken())

	{
		organisationGroup.POST("/register", middlewares.ValidateToken(), organisations.RegisterOrganisation) // Register a new organization
		organisationGroup.GET("/get/:id", organisations.GetOrganisation)
		organisationGroup.PUT("/update/:id", organisations.UpdateOrganisation)
	}

	// Employee management routes
	employeeGroup := r.Group("/api/v1/employees")
	employeeGroup.Use(middlewares.ValidateToken())
	{
		employeeGroup.GET("/is-employee/:user_id", employees.IsEmployee)
		employeeGroup.POST("/add", employees.AddEmployees)
		employeeGroup.GET("/get/:organisation_id", employees.GetEmployees)
		employeeGroup.GET("/list/:id", employees.GetEmployeeByID)
		employeeGroup.PUT("/update/:id", employees.UpdateEmployee)
		employeeGroup.DELETE("/delete/:id", employees.DeleteEmployee)
		employeeGroup.GET("/search", employees.SearchEmployee)
		employeeGroup.POST("/make-manager", employees.MakeManager)
	}
	// Task management routes
	taskGroup := r.Group("/api/v1/tasks")
	taskGroup.Use(middlewares.ValidateToken())
	{
		taskGroup.POST("/create", tasks.CreateTask)
		taskGroup.PUT("/update/:id", tasks.UpdateTask)
		taskGroup.DELETE("/delete/:id", tasks.SoftDeleteTask)
		taskGroup.GET("/get/:id", tasks.GetTaskByID)
		taskGroup.GET("/team", tasks.GetTeamTasks)
		taskGroup.GET("/project", tasks.GetProjectTasks)
		taskGroup.GET("/personal", tasks.GetPersonalTasks)
		taskGroup.GET("/organisation/:organisation_id", tasks.GetOrganisationTasks)
		taskGroup.GET("/title/:title", tasks.GetTaskByTitle)
	}

	{
		r.POST("/api/v1/create/boss", middlewares.ValidateToken(), employees.CreateBoss)
		r.GET("/api/v1/get/boss", middlewares.ValidateToken(), employees.GetBoss)
		r.POST("/api/v1/create/employee", middlewares.ValidateToken(), employees.CreateEmployee)
	}
	// Attendance routes
	attendanceGroup := r.Group("/api/v1/attendence")
	attendanceGroup.Use(middlewares.ValidateToken())
	{
		attendanceGroup.POST("/create", attendence.PostAttendance)
		attendanceGroup.PUT("/:id", attendence.UpdateAttendance)
		attendanceGroup.GET("/:id", attendence.GetAttendanceByDateAndUser)
	}
	// Team management routes
	teamGroup := r.Group("/api/v1/teams")
	teamGroup.Use(middlewares.ValidateToken())
	{
		teamGroup.POST("/create", teams.CreateTeam)
		teamGroup.GET("/:id", teams.GetTeam)
		teamGroup.PUT("/:id", teams.UpdateTeam)
		teamGroup.GET("/organisation/:organisation_id", teams.GetTeamsByOrganisation)
		teamGroup.DELETE("/delete/:id", teams.DelteTeam)
	}
	// Project management routes
	projectGroup := r.Group("/api/v1/projects")
	projectGroup.Use(middlewares.ValidateToken())
	{
		projectGroup.POST("/create", projects.CreateProject)
		projectGroup.GET("/:id", projects.GetProject)
		projectGroup.PUT("/:id", projects.UpdateProject)
		projectGroup.GET("/organisation/:organisation_id", projects.GetProjectsByOrganisation)
		projectGroup.GET("/team/:team_id", projects.GetProjectsByTeam)
		projectGroup.GET("/less-data/:organisation_id", projects.GetProjectsWithLessData)

	}
	//calls
	r.GET("/api/v1/call/:id", calls.HandleRoom)
}
