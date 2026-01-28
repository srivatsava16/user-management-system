package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Prevent MIME type sniffing
		ctx.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		ctx.Header("X-Frame-Options", "DENY")

		// XSS Protection (legacy browsers)
		ctx.Header("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		ctx.Header("Content-Security-Policy", "default-src 'self'")

		// Strict Transport Security (enable when using HTTPS)
		// ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		ctx.Next()
	}
}
