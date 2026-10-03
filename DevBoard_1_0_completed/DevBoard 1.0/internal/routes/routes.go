package routes

import (
	"github.com/Demiladeolorunsola/devboard-1.0/internal/config"
	"github.com/Demiladeolorunsola/devboard-1.0/internal/handlers"
	"github.com/Demiladeolorunsola/devboard-1.0/internal/middleware"
	"github.com/Demiladeolorunsola/devboard-1.0/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(
	router *gin.Engine,
	db *gorm.DB,
	config *config.Config,
) {
	authService := services.NewAuthService(db)
	userService := services.NewUserService(db)
	projectService := services.NewProjectService(db)
	taskService := services.NewTaskService(db)
	timeService := services.NewTimeService(db)
	reportService := services.NewReportService(db)

	authHandler := handlers.NewAuthHandler(authService, config.JWTSecret)
	userHandler := handlers.NewUserHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	timeHandler := handlers.NewTimeHandler(timeService)
	reportHandler := handlers.NewReportHandler(reportService)

	// =========================
	// Auth (public)
	// =========================

	auth := router.Group("/api/auth")

	auth.POST(
		"/register",
		authHandler.Register,
	)

	auth.POST(
		"/login",
		authHandler.Login,
	)

	protected := router.Group("/api")
	protected.Use(
		middleware.AuthMiddleware(config.JWTSecret),
	)

	// =========================
	// Profile
	// =========================

	protected.GET(
		"/profile",
		userHandler.Profile,
	)

	protected.PUT(
		"/profile",
		userHandler.UpdateProfile,
	)

	protected.DELETE(
		"/profile",
		userHandler.DeleteProfile,
	)

	// =========================
	// Projects
	// =========================

	protected.POST(
		"/projects",
		projectHandler.Create,
	)

	protected.GET(
		"/projects",
		projectHandler.GetAll,
	)

	protected.GET(
		"/projects/:project_id",
		projectHandler.GetByID,
	)

	protected.PUT(
		"/projects/:project_id",
		projectHandler.Update,
	)

	protected.DELETE(
		"/projects/:project_id",
		projectHandler.Delete,
	)

	// =========================
	// Tasks
	// =========================

	protected.POST(
		"/projects/:project_id/tasks",
		taskHandler.Create,
	)

	protected.GET(
		"/projects/:project_id/tasks",
		taskHandler.GetByProject,
	)

	protected.GET(
		"/tasks/:task_id",
		taskHandler.GetByID,
	)

	protected.PUT(
		"/tasks/:task_id",
		taskHandler.Update,
	)

	protected.DELETE(
		"/tasks/:task_id",
		taskHandler.Delete,
	)

	// =========================
	// Time Tracking
	// =========================

	protected.POST(
		"/time/start",
		timeHandler.Start,
	)

	protected.POST(
		"/time/stop/:id",
		timeHandler.Stop,
	)

	protected.GET(
		"/time",
		timeHandler.GetAll,
	)

	protected.GET(
		"/projects/:project_id/time",
		timeHandler.GetByProject,
	)

	protected.GET(
		"/tasks/:task_id/time",
		timeHandler.GetByTask,
	)

	// =========================
	// Reports
	// =========================

	// Project time summary
	protected.GET(
		"/reports/projects/:project_id/time",
		reportHandler.ProjectTimeSummary,
	)

	// Task productivity summary
	protected.GET(
		"/reports/tasks/productivity",
		reportHandler.TaskProductivity,
	)

	// Overall user activity
	protected.GET(
		"/reports/activity",
		reportHandler.UserActivity,
	)

	// Overall dashboard report
	protected.GET(
		"/reports/dashboard",
		reportHandler.Dashboard,
	)
}
