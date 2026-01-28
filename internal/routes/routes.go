package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"user-management-system/internal/constants"
	"user-management-system/internal/controllers"
	"user-management-system/internal/middleware"
	"user-management-system/internal/repository"
	"user-management-system/internal/services"
	"user-management-system/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// CORS Configuration
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	corsConfig := cors.Config{
		AllowOrigins:     strings.Split(frontendURL, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// Security Headers
	router.Use(middleware.SecurityHeaders())

	// Health Check Endpoint (no auth required)
	router.GET("/health", func(ctx *gin.Context) {
		// Check database connection
		db := database.GetDB()
		sqlDB, err := db.DB()
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"database": "error",
				"error":    err.Error(),
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"database": "down",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"database": "connected",
		})
	})

	db := database.GetDB()

	auditRepo := repository.NewAuditRepo(db)
	auditService := services.NewAuditService(auditRepo)
	auditController := controllers.NewAuditController(auditService)

	tokenInvalidationRepo := repository.NewTokenInvalidationRepo(db)
	tokenInvalidationService := services.NewTokenInvalidationService(tokenInvalidationRepo)

	buRepo := repository.NewBusinessUnitRepo(db)
	buService := services.NewBusinessUnitService(buRepo)
	buController := controllers.NewBusinessUnitController(buService, auditService)

	divRepo := repository.NewDivisionRepo(db)
	divService := services.NewDivisionService(divRepo)
	divController := controllers.NewDivisionController(divService, auditService)

	userRepo := repository.NewUserRepo(db)
	userService := services.NewUserService(userRepo, tokenInvalidationService, auditService)
	userController := controllers.NewUserController(userService, auditService, tokenInvalidationService)

	roleRepo := repository.NewRoleRepo(db)
	roleService := services.NewRoleService(roleRepo, userRepo)
	roleController := controllers.NewRoleController(roleService, auditService, tokenInvalidationService, userRepo)

	jwtService, err := services.NewJWTService()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize JWT Service: %v", err))
	}
	authController := controllers.NewAuthController(userService, roleService, jwtService, tokenInvalidationService)

	approvalService := services.NewApprovalService(userRepo, divRepo, buRepo)
	approvalController := controllers.NewApprovalController(approvalService, auditService)

	// API v1 routes
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", middleware.LoginRateLimiter(), authController.Login)
		auth.POST("/forgot-password", middleware.PasswordResetRateLimiter(), authController.ForgotPassword)
		auth.POST("/reset-password/:token", authController.ResetPassword)
	}

	api := router.Group("/api/v1")

	api.Use(middleware.RequireJWT(jwtService, tokenInvalidationService))
	{
		api.POST("/logout", authController.Logout)

		approvalRoutes := api.Group("/approvals")

		approvalRoutes.Use(middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin, constants.RoleDVAdmin))
		{
			approvalRoutes.PUT("/approve", approvalController.ApproveEntity)
			approvalRoutes.GET("/pending/:entity_type", approvalController.GetPendingEntities)
		}

		buRoutes := api.Group("/business-units")
		buRoutes.Use(middleware.RequireRole(constants.RoleSuperAdmin))
		{
			buRoutes.POST("", buController.CreateBusinessUnit)
			buRoutes.PUT("/:id", buController.UpdateBusinessUnit)
		}
		api.GET("/business-units/:id", middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin, constants.RoleDVAdmin), buController.GetBusinessUnitByID)
		api.GET("/business-units", middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin, constants.RoleDVAdmin), buController.GetAllBusinessUnits)

		divRoutes := api.Group("/divisions")
		divRoutes.Use(middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin))
		{
			divRoutes.POST("", divController.CreateDivision)
			divRoutes.PUT("/:id", divController.UpdateDivision)
		}

		api.GET("/divisions/:id", middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin, constants.RoleDVAdmin), divController.GetDivisionByID)
		api.GET("/divisions", middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin, constants.RoleDVAdmin), divController.GetAllDivisions)
		api.GET("/business-units/:id/divisions", middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin, constants.RoleDVAdmin), divController.GetDivisionsByBusinessUnit)

		roleRoutes := api.Group("/roles")
		roleRoutes.Use(middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin, constants.RoleDVAdmin))
		{
			roleRoutes.POST("", roleController.CreateRole)
			roleRoutes.PUT("/:id", roleController.UpdateRole)
			roleRoutes.GET("/:id", roleController.GetRoleByID)
			roleRoutes.GET("", roleController.GetAllRoles)
			roleRoutes.GET("/users/:id/permissions", roleController.GetUserPermissions)
		}

		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.RequireRole(constants.RoleSuperAdmin, constants.RoleBUAdmin, constants.RoleDVAdmin))
		{
			userRoutes.POST("", userController.CreateUser)
			userRoutes.GET("/:id", userController.GetUserByID)
			userRoutes.GET("", userController.GetAllUsers)
			userRoutes.PUT("/:id", userController.UpdateUser)
			userRoutes.GET("/divisions/:id/users", userController.GetUsersByDivision)
			userRoutes.GET("/business-units/:id/users", userController.GetUsersByBusinessUnit)
		}

		api.GET("/audit-logs", auditController.GetAllAuditLogs)

	}

	return router
}
