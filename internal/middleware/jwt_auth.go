package middleware

import (
	"errors"
	"net/http"
	"strings"
	"user-management-system/internal/services"
	"user-management-system/internal/shared"

	"github.com/gin-gonic/gin"
)

func RequireJWT(jwtService *services.JWTService, tokenInvalidationService *services.TokenInvalidationService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format. use: Bearer <token>"})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtService.ValidateToken(token)

		if err != nil {
			var errorMessage string
			switch {
			case errors.Is(err, shared.ErrTokenExpired):
				errorMessage = "Token has expired. Please log in again."
			case errors.Is(err, shared.ErrTokenInvalid):
				errorMessage = "Invalid token. Please log in again."
			case errors.Is(err, shared.ErrTokenMalformed):
				errorMessage = "Malformed token. Please log in again."
			case errors.Is(err, shared.ErrTokenInvalidSignature):
				errorMessage = "Invalid token signature. Please log in again."
			default:
				errorMessage = "Unauthorized access. Please log in."
			}

			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
			ctx.Abort()
			return
		}
		if tokenInvalidationService.IsTokenInvalidated(claims.JTI) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user has been logged out. Please log in again."})
			ctx.Abort()
			return
		}

		if tokenInvalidationService.IsUserTokenInvalidated(claims.UserID) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Your token is no longer valid due to recent account changes. Please log in again."})
			ctx.Abort()
			return
		}
		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_roles", claims.Roles)
		ctx.Set("business_unit_id", claims.BusinessUnitID)
		ctx.Set("division_id", claims.DivisionID)
		ctx.Set("jti", claims.JTI)

		ctx.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userRoles, exists := ctx.Get("user_roles")
		if !exists {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			ctx.Abort()
			return
		}

		hasPermission := checkUserRoles(userRoles, allowedRoles)

		if !hasPermission {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			ctx.Abort()
			return
		}
		ctx.Next()

	}
}

func checkUserRoles(userRoles interface{}, allowedRoles []string) bool {

	roleNames, ok := userRoles.([]string)
	if !ok {
		return false
	}

	for _, userRole := range roleNames {
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return true
			}
		}
	}

	return false
}
