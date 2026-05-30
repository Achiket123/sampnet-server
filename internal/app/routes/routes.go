package routes

import (
	attendanceHTTP "server/internal/adapters/http/attendance"
	authHTTP "server/internal/adapters/http/auth"
	callHTTP "server/internal/adapters/http/calls"
	callStateHTTP "server/internal/adapters/http/callstate"
	chatHTTP "server/internal/adapters/http/chats"
	employeeHTTP "server/internal/adapters/http/employees"
	fileHTTP "server/internal/adapters/http/files"
	messageHTTP "server/internal/adapters/http/messages"
	organisationHTTP "server/internal/adapters/http/organisation"
	projectHTTP "server/internal/adapters/http/projects"
	taskHTTP "server/internal/adapters/http/tasks"
	teamHTTP "server/internal/adapters/http/teams"
	attendanceRepo "server/internal/adapters/repository/attendance"
	authRepo "server/internal/adapters/repository/auth"
	callStateRepo "server/internal/adapters/repository/callstate"
	chatRepo "server/internal/adapters/repository/chats"
	employeeRepo "server/internal/adapters/repository/employees"
	fileRepo "server/internal/adapters/repository/files"
	messageRepo "server/internal/adapters/repository/messages"
	organisationRepo "server/internal/adapters/repository/organisation"
	projectRepo "server/internal/adapters/repository/projects"
	taskRepo "server/internal/adapters/repository/tasks"
	teamRepo "server/internal/adapters/repository/teams"
	"server/internal/platform/database"
	"server/internal/platform/middlewares"
	attendanceService "server/internal/usecase/attendance"
	authService "server/internal/usecase/auth"
	callService "server/internal/usecase/calls"
	callStateService "server/internal/usecase/callstate"
	chatService "server/internal/usecase/chats"
	employeeService "server/internal/usecase/employees"
	fileService "server/internal/usecase/files"
	messageService "server/internal/usecase/messages"
	organisationUsecase "server/internal/usecase/organisation"
	projectService "server/internal/usecase/projects"
	taskService "server/internal/usecase/tasks"
	teamService "server/internal/usecase/teams"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	validateToken := middlewares.ValidateToken()

	organisationRepository := organisationRepo.NewGormRepository(database.DB)
	organisationService := organisationUsecase.NewService(organisationRepository)
	organisationHandler := organisationHTTP.NewHandler(organisationService)
	organisationHTTP.RegisterRoutes(r, organisationHandler, validateToken)

	authRepository := authRepo.NewGormRepository(database.DB)
	empRepository := employeeRepo.NewGormRepository(database.DB)
	authUseCase := authService.NewService(authRepository, empRepository)
	authHandler := authHTTP.NewHandler(authUseCase)
	authHTTP.RegisterRoutes(r, authHandler, validateToken)

	employeeUseCase := employeeService.NewService(empRepository, authRepository, organisationRepository)
	employeeHandler := employeeHTTP.NewHandler(employeeUseCase)
	employeeHTTP.RegisterRoutes(r, employeeHandler, validateToken)

	taskRepository := taskRepo.NewGormRepository(database.DB)
	taskUseCase := taskService.NewService(taskRepository)
	taskHandler := taskHTTP.NewHandler(taskUseCase)
	taskHTTP.RegisterRoutes(r, taskHandler, validateToken)

	attendanceRepository := attendanceRepo.NewGormRepository(database.DB)
	attendanceUseCase := attendanceService.NewService(attendanceRepository)
	attendanceHandler := attendanceHTTP.NewHandler(attendanceUseCase)
	attendanceHTTP.RegisterRoutes(r, attendanceHandler, validateToken)

	teamRepository := teamRepo.NewGormRepository(database.DB)
	teamUseCase := teamService.NewService(teamRepository)
	teamHandler := teamHTTP.NewHandler(teamUseCase)
	teamHTTP.RegisterRoutes(r, teamHandler, validateToken)

	projectRepository := projectRepo.NewGormRepository(database.DB)
	projectUseCase := projectService.NewService(projectRepository)
	projectHandler := projectHTTP.NewHandler(projectUseCase)
	projectHTTP.RegisterRoutes(r, projectHandler, validateToken)

	fileRepository := fileRepo.NewGormRepository(database.DB)
	fileUseCase := fileService.NewService(fileRepository)
	fileHandler := fileHTTP.NewHandler(fileUseCase)
	fileHTTP.RegisterRoutes(r, fileHandler)

	chatRepository := chatRepo.NewGormRepository(database.DB)
	chatUseCase := chatService.NewService(chatRepository)
	chatHandler := chatHTTP.NewHandler(chatUseCase)
	chatHTTP.RegisterRoutes(r, chatHandler, validateToken)

	messageRepository := messageRepo.NewGormRepository(database.DB)
	messageUseCase := messageService.NewService(messageRepository, chatRepository)
	messageHandler := messageHTTP.NewHandler(messageUseCase)
	messageHTTP.RegisterRoutes(r, messageHandler, validateToken)

	callStateRepository := callStateRepo.NewGormRepository(database.DB)
	callStateUseCase := callStateService.NewService(callStateRepository)
	callStateHandler := callStateHTTP.NewHandler(callStateUseCase)
	callStateHTTP.RegisterRoutes(r, callStateHandler, validateToken)

	callUseCase := callService.NewService()
	callHandler := callHTTP.NewHandler(callUseCase)
	callHTTP.RegisterRoutes(r, callHandler)
}
