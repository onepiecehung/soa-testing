package middleware

import (
	"bytes"
	"io"
	"time"

	"product-management/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AutoLogger middleware automatically logs request and response
func AutoLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Log request
		requestLogger := logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// Log request body if exists
		if c.Request.Body != nil {
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			if len(body) > 0 {
				requestLogger = requestLogger.WithField("request_body", string(body))
			}
		}

		requestLogger.Info("Incoming request")

		// Create a custom response writer to capture response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Stop timer
		duration := time.Since(start)

		// Log response
		responseLogger := logger.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"duration": duration,
		})

		// Log response body if exists
		if blw.body.Len() > 0 {
			responseLogger = responseLogger.WithField("response_body", blw.body.String())
		}

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				responseLogger = responseLogger.WithField("error", err.Error())
			}
			responseLogger.Error("Request completed with errors")
		} else {
			responseLogger.Info("Request completed successfully")
		}
	}
}

// bodyLogWriter is a custom response writer to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
