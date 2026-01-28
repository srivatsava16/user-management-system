package routes

import (
	"fmt"
	"user-management-system/internal/controllers"
	"user-management-system/internal/middleware"
	"user-management-system/internal/repository"
	"user-management-system/internal/services"
	"user-management-system/pkg/database"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

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

	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.POST("/reset-password/:token", authController.ResetPassword)

	}

	api := router.Group("/api")

	api.Use(middleware.RequireJWT(jwtService, tokenInvalidationService))
	{
		api.POST("/logout", authController.Logout)

		approvalRoutes := api.Group("/approvals")

		approvalRoutes.Use(middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"))
		{
			approvalRoutes.PUT("/approve", approvalController.ApproveEntity)
			approvalRoutes.GET("/pending/:entity_type", approvalController.GetPendingEntities)
		}

		buRoutes := api.Group("/business-units")
		buRoutes.Use(middleware.RequireRole("Super Admin"))
		{
			buRoutes.POST("", buController.CreateBusinessUnit)
			buRoutes.PUT("/:id", buController.UpdateBusinessUnit)
		}
		api.GET("/business-units/:id", middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"), buController.GetBusinessUnitByID)
		api.GET("/business-units", middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"), buController.GetAllBusinessUnits)

		divRoutes := api.Group("/divisions")
		divRoutes.Use(middleware.RequireRole("Super Admin", "BU Admin"))
		{
			divRoutes.POST("", divController.CreateDivision)
			divRoutes.PUT("/:id", divController.UpdateDivision)
		}

		api.GET("/divisions/:id", middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"), divController.GetDivisionByID)
		api.GET("/divisions", middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"), divController.GetAllDivisions)
		api.GET("/business-units/:id/divisions", middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"), divController.GetDivisionsByBusinessUnit)

		roleRoutes := api.Group("/roles")
		roleRoutes.Use(middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"))
		{
			roleRoutes.POST("", roleController.CreateRole)
			roleRoutes.PUT("/:id", roleController.UpdateRole)
			roleRoutes.GET("/:id", roleController.GetRoleByID)
			roleRoutes.GET("", roleController.GetAllRoles)
			roleRoutes.GET("/users/:id/permissions", roleController.GetUserPermissions)
		}

		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.RequireRole("Super Admin", "BU Admin", "DV Admin"))
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
