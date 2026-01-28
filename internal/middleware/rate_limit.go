package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     time.Duration
	burst    int
}

type visitor struct {
	lastSeen time.Time
	tokens   int
}

func newRateLimiter(rate time.Duration, burst int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}

	// Cleanup old visitors every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 10*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

func (rl *rateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			lastSeen: time.Now(),
			tokens:   rl.burst,
		}
		rl.visitors[ip] = v
		return v
	}

	// Refill tokens based on time passed
	elapsed := time.Since(v.lastSeen)
	tokensToAdd := int(elapsed / rl.rate)
	if tokensToAdd > 0 {
		v.tokens += tokensToAdd
		if v.tokens > rl.burst {
			v.tokens = rl.burst
		}
		v.lastSeen = time.Now()
	}

	return v
}

func (rl *rateLimiter) allow(ip string) bool {
	v := rl.getVisitor(ip)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware creates a rate limiting middleware
// rate: duration between each allowed request
// burst: maximum number of requests allowed in burst
func RateLimitMiddleware(rate time.Duration, burst int) gin.HandlerFunc {
	limiter := newRateLimiter(rate, burst)

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		if !limiter.allow(ip) {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// LoginRateLimiter limits login attempts to 5 per 15 minutes per IP
func LoginRateLimiter() gin.HandlerFunc {
	return RateLimitMiddleware(3*time.Minute, 5)
}

// PasswordResetRateLimiter limits password reset requests to 3 per hour per IP
func PasswordResetRateLimiter() gin.HandlerFunc {
	return RateLimitMiddleware(20*time.Minute, 3)
}
