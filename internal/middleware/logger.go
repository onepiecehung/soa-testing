package middleware

import (
	"time"

	"product-management/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RequestLogger middleware logs all requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		duration := time.Since(start)

		// Log request details
		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("Request completed")
	}
}

// ErrorLogger middleware logs all errors
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.WithFields(logrus.Fields{
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
					"status": c.Writer.Status(),
					"error":  err.Error(),
				}).Error("Request error")
			}
		}
	}
}

// Recovery middleware logs panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.WithFields(logrus.Fields{
					"method": c.Request.Method,
					"path":   c.Request.URL.Path,
					"error":  err,
				}).Error("Panic recovered")
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
