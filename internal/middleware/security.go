package middleware

import (
	"crypto/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// XSSMiddleware provides protection against XSS attacks
func XSSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set security headers
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Sanitize input for POST and PUT requests
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			// TODO: Implement input sanitization logic here
		}

		c.Next()
	}
}

// CSRFMiddleware provides protection against CSRF attacks
func CSRFMiddleware() gin.HandlerFunc {
	// Generate a secure random key for CSRF protection
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic("Failed to generate CSRF key: " + err.Error())
	}

	// Create CSRF middleware with secure settings
	csrfMiddleware := csrf.Protect(
		key,
		csrf.Secure(true),                      // Only send cookies over HTTPS
		csrf.HttpOnly(true),                    // Prevent JavaScript access to cookies
		csrf.MaxAge(3600),                      // Token expires after 1 hour
		csrf.Path("/"),                         // Cookie path
		csrf.SameSite(csrf.SameSiteStrictMode), // Strict same-site policy
	)

	return func(c *gin.Context) {
		// Convert Gin context to http.Handler
		csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Writer = w.(gin.ResponseWriter)
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}

// RateLimitMiddleware limits the number of requests from a single IP
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	// Create a map to store request counts
	requests := make(map[string]int)
	lastReset := make(map[string]time.Time)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		// Reset counter if window has passed
		if lastResetTime, exists := lastReset[ip]; exists {
			if now.Sub(lastResetTime) > window {
				requests[ip] = 0
				lastReset[ip] = now
			}
		} else {
			lastReset[ip] = now
		}

		// Check if limit has been reached
		if requests[ip] >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		// Increment request count
		requests[ip]++

		c.Next()
	}
}
