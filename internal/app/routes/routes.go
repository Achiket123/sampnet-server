package routes

import (
	analyticsHTTP "server/internal/adapters/http/analytics"
	attendanceHTTP "server/internal/adapters/http/attendance"
	authHTTP "server/internal/adapters/http/auth"
	callHTTP "server/internal/adapters/http/calls"
	callStateHTTP "server/internal/adapters/http/callstate"
	chatHTTP "server/internal/adapters/http/chats"
	employeeHTTP "server/internal/adapters/http/employees"
	inviteHTTP "server/internal/adapters/http/invites"
	fileHTTP "server/internal/adapters/http/files"
	leaveHTTP "server/internal/adapters/http/leave"
	messageHTTP "server/internal/adapters/http/messages"
	notificationsHTTP "server/internal/adapters/http/notifications"
	onboardingHTTP "server/internal/adapters/http/onboarding"
	organisationHTTP "server/internal/adapters/http/organisation"
	projectHTTP "server/internal/adapters/http/projects"
	taskHTTP "server/internal/adapters/http/tasks"
	taskAttachmentHTTP "server/internal/adapters/http/task_attachments"
	taskCommentHTTP "server/internal/adapters/http/task_comments"
	searchHTTP "server/internal/adapters/http/search"
	teamHTTP "server/internal/adapters/http/teams"
	resourcesHTTP "server/internal/adapters/http/resources"
	calendarHTTP "server/internal/adapters/http/calendar"
	researchHTTP "server/internal/adapters/http/research"
	peopleHTTP "server/internal/adapters/http/people"
	settingsHTTP "server/internal/adapters/http/settings"
	analyticsRepo "server/internal/adapters/repository/analytics"
	attendanceRepo "server/internal/adapters/repository/attendance"
	authRepo "server/internal/adapters/repository/auth"
	callStateRepo "server/internal/adapters/repository/callstate"
	chatRepo "server/internal/adapters/repository/chats"
	employeeRepo "server/internal/adapters/repository/employees"
	inviteRepo "server/internal/adapters/repository/invites"
	fileRepo "server/internal/adapters/repository/files"
	leaveRepo "server/internal/adapters/repository/leave"
	messageRepo "server/internal/adapters/repository/messages"
	notificationsRepo "server/internal/adapters/repository/notifications"
	onboardingRepo "server/internal/adapters/repository/onboarding"
	organisationRepo "server/internal/adapters/repository/organisation"
	projectRepo "server/internal/adapters/repository/projects"
	taskRepo "server/internal/adapters/repository/tasks"
	taskAttachmentRepo "server/internal/adapters/repository/task_attachments"
	taskCommentRepo "server/internal/adapters/repository/task_comments"
	searchRepo "server/internal/adapters/repository/search"
	teamRepo "server/internal/adapters/repository/teams"
	resourcesRepo "server/internal/adapters/repository/resources"
	calendarRepo "server/internal/adapters/repository/calendar"
	researchRepo "server/internal/adapters/repository/research"
	peopleRepo "server/internal/adapters/repository/people"
	settingsRepo "server/internal/adapters/repository/settings"
	"server/internal/platform/database"
	"server/internal/platform/middlewares"
	mailerPlatform "server/internal/platform/mailer"
	"server/internal/platform/websocket"
	analyticsService "server/internal/usecase/analytics"
	attendanceService "server/internal/usecase/attendance"
	authService "server/internal/usecase/auth"
	callService "server/internal/usecase/calls"
	callStateService "server/internal/usecase/callstate"
	chatService "server/internal/usecase/chats"
	employeeService "server/internal/usecase/employees"
	inviteService "server/internal/usecase/invites"
	fileService "server/internal/usecase/files"
	leaveService "server/internal/usecase/leave"
	messageService "server/internal/usecase/messages"
	notificationsService "server/internal/usecase/notifications"
	onboardingService "server/internal/usecase/onboarding"
	organisationUsecase "server/internal/usecase/organisation"
	projectService "server/internal/usecase/projects"
	taskService "server/internal/usecase/tasks"
	taskAttachmentService "server/internal/usecase/task_attachments"
	taskCommentService "server/internal/usecase/task_comments"
	searchService "server/internal/usecase/search"
	teamService "server/internal/usecase/teams"
	resourcesService "server/internal/usecase/resources"
	calendarService "server/internal/usecase/calendar"
	researchService "server/internal/usecase/research"
	peopleService "server/internal/usecase/people"
	settingsService "server/internal/usecase/settings"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	validateToken := middlewares.ValidateToken()

	// Initialize WebSocket hub and broadcaster
	hub := websocket.NewHub()
	go hub.Run()
	broadcaster := websocket.NewBroadcaster(hub)
	wsManager := websocket.NewManager(hub)

	notificationRepository := notificationsRepo.NewGormRepository(database.DB)
	notificationUseCase := notificationsService.NewService(notificationRepository, broadcaster)
	notificationHandler := notificationsHTTP.NewHandler(notificationUseCase, hub)
	notificationsHTTP.RegisterRoutes(r, notificationHandler, validateToken)

	leaveRepository := leaveRepo.NewGormRepository(database.DB)
	leaveUseCase := leaveService.NewService(leaveRepository, notificationUseCase)
	leaveHandler := leaveHTTP.NewHandler(leaveUseCase)
	leaveHTTP.RegisterRoutes(r, leaveHandler, validateToken)

	onboardingRepository := onboardingRepo.NewGormRepository(database.DB)
	onboardingUseCase := onboardingService.NewService(onboardingRepository)
	onboardingHandler := onboardingHTTP.NewHandler(onboardingUseCase)
	onboardingHTTP.RegisterRoutes(r, onboardingHandler, validateToken)

	organisationRepository := organisationRepo.NewGormRepository(database.DB)
	organisationService := organisationUsecase.NewService(organisationRepository)
	organisationHandler := organisationHTTP.NewHandler(organisationService)
	organisationHTTP.RegisterRoutes(r, organisationHandler, validateToken)

	gmailMailer := mailerPlatform.NewGmailMailer()
	authRepository := authRepo.NewGormRepository(database.DB)
	empRepository := employeeRepo.NewGormRepository(database.DB)
	authUseCase := authService.NewService(authRepository, empRepository, gmailMailer)
	authHandler := authHTTP.NewHandler(authUseCase)
	authHTTP.RegisterRoutes(r, authHandler, validateToken)

	employeeUseCase := employeeService.NewService(empRepository, authRepository, organisationRepository)
	employeeHandler := employeeHTTP.NewHandler(employeeUseCase)
	employeeHTTP.RegisterRoutes(r, employeeHandler, validateToken)

	inviteRepository := inviteRepo.NewGormRepository(database.DB)
	inviteUseCase := inviteService.NewService(inviteRepository, authRepository, empRepository, gmailMailer)
	inviteHandler := inviteHTTP.NewHandler(inviteUseCase)
	inviteHTTP.RegisterRoutes(r, inviteHandler, validateToken)

	taskRepository := taskRepo.NewGormRepository(database.DB)
	taskUseCase := taskService.NewService(taskRepository, notificationUseCase)
	taskHandler := taskHTTP.NewHandler(taskUseCase)
	taskHTTP.RegisterRoutes(r, taskHandler, validateToken)

	taskCommentRepository := taskCommentRepo.NewGormRepository(database.DB)
	taskCommentUseCase := taskCommentService.NewService(taskCommentRepository, taskRepository, notificationUseCase)
	taskCommentHandler := taskCommentHTTP.NewHandler(taskCommentUseCase)
	taskCommentHTTP.RegisterRoutes(r, taskCommentHandler, validateToken)

	taskAttachmentRepository := taskAttachmentRepo.NewGormRepository(database.DB)
	taskAttachmentUseCase := taskAttachmentService.NewService(taskAttachmentRepository)
	taskAttachmentHandler := taskAttachmentHTTP.NewHandler(taskAttachmentUseCase)
	taskAttachmentHTTP.RegisterRoutes(r, taskAttachmentHandler, validateToken)

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
	messageUseCase := messageService.NewService(messageRepository, chatRepository, wsManager)
	messageHandler := messageHTTP.NewHandler(messageUseCase)
	messageHTTP.RegisterRoutes(r, messageHandler, validateToken)

	callStateRepository := callStateRepo.NewGormRepository(database.DB)
	callStateUseCase := callStateService.NewService(callStateRepository)
	callStateHandler := callStateHTTP.NewHandler(callStateUseCase, wsManager, chatRepository)
	callStateHTTP.RegisterRoutes(r, callStateHandler, validateToken)

	callUseCase := callService.NewService()
	callHandler := callHTTP.NewHandler(callUseCase)
	callHTTP.RegisterRoutes(r, callHandler)

	searchRepository := searchRepo.NewGormRepository(database.DB)
	searchUseCase := searchService.NewService(searchRepository)
	searchHandler := searchHTTP.NewHandler(searchUseCase)
	searchHTTP.RegisterRoutes(r, searchHandler, validateToken)

	resourcesRepository := resourcesRepo.NewGormRepository(database.DB)
	resourcesUseCase := resourcesService.NewService(resourcesRepository)
	resourcesHandler := resourcesHTTP.NewHandler(resourcesUseCase)
	resourcesHTTP.RegisterRoutes(r, resourcesHandler, validateToken)

	calendarRepository := calendarRepo.NewGormRepository(database.DB)
	calendarUseCase := calendarService.NewService(calendarRepository)
	calendarHandler := calendarHTTP.NewHandler(calendarUseCase)
	calendarHTTP.RegisterRoutes(r, calendarHandler, validateToken)

	researchRepository := researchRepo.NewGormRepository(database.DB)
	researchUseCase := researchService.NewUseCase(researchRepository, projectRepository, teamRepository, notificationUseCase)
	researchHandler := researchHTTP.NewHandler(researchUseCase)
	researchHTTP.RegisterRoutes(r, researchHandler, validateToken)

	peopleRepository := peopleRepo.NewRepository(database.DB)
	peopleUseCase := peopleService.NewUseCase(peopleRepository)
	peopleHandler := peopleHTTP.NewHandler(peopleUseCase)
	peopleHTTP.RegisterRoutes(r, peopleHandler, validateToken)

	analyticsRepository := analyticsRepo.NewGormRepository(database.DB)
	analyticsUseCase := analyticsService.NewService(analyticsRepository)
	analyticsHandler := analyticsHTTP.NewHandler(analyticsUseCase)
	
	// Ensure Boss or Manager role for analytics
	roleMiddleware := middlewares.RoleMiddleware([]string{"boss", "manager"})
	analyticsHTTP.RegisterRoutes(r, analyticsHandler, validateToken, roleMiddleware)

	settingsRepository := settingsRepo.NewGormRepository(database.DB)
	settingsUseCase := settingsService.NewUseCase(settingsRepository)
	settingsHandler := settingsHTTP.NewHandler(settingsUseCase)
	settingsHTTP.RegisterRoutes(r, settingsHandler, validateToken)
}
