package controllers

import (
	"net/http"
	"user-management-system/internal/services"
	"user-management-system/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	userService              *services.UserService
	roleService              *services.RoleService
	jwtService               *services.JWTService
	tokenInvalidationService *services.TokenInvalidationService
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func NewAuthController(userService *services.UserService, roleService *services.RoleService, jwtService *services.JWTService, tokenInvalidationService *services.TokenInvalidationService) *AuthController {
	return &AuthController{userService: userService, roleService: roleService, jwtService: jwtService, tokenInvalidationService: tokenInvalidationService}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := c.userService.ValidateLogin(req.Email, req.Password)

	if err != nil {
		log.Warn().Str("email", req.Email).Msg("Failed login attempt")

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	roleNames, err := c.roleService.GetRoleNamesByIDs(user.RoleIds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user roles"})
		return
	}

	err = c.tokenInvalidationService.ClearUserInvalidations(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear token invalidations"})
		return
	}

	err = c.tokenInvalidationService.CleanupExpiredInvalidations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup expired token invalidations"})
		return
	}

	token, err := c.jwtService.GenerateToken(user.ID, user.Email, roleNames, user.BusinessUnitID, user.DivisionID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "roles": roleNames, "user_id": user.ID, "token": token})
}

func (c *AuthController) Logout(ctx *gin.Context) {

	jti, exists := ctx.Get("jti")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to retrieve token identifier"})
		return
	}

	userID := shared.GetUserIDFromContext(ctx)

	err := c.tokenInvalidationService.InvalidateToken(jti.(string), userID, "logout")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	err := c.userService.ForgotPassword(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to email"})
}

func (c *AuthController) ResetPassword(ctx *gin.Context) {

	token := ctx.Param("token")

	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	userID, err := c.userService.ResetPassword(token, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.tokenInvalidationService.InvalidateUserTokens(userID, "password_reset")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate tokens"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password Updated Successfully"})
}
