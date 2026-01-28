package shared

import "github.com/gin-gonic/gin"

func GetUserFromContext(ctx *gin.Context) string {
	if userEmail, exists := ctx.Get("user_email"); exists && userEmail != nil {
		return userEmail.(string)
	}
	return "system"
}

func GetUserIDFromContext(ctx *gin.Context) int {
	if userID, exists := ctx.Get("user_id"); exists && userID != nil {
		return userID.(int)
	}
	return 0
}
