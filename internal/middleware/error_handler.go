package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			status := c.Writer.Status()
			if status < 400 {
				status = http.StatusInternalServerError
			}

			c.JSON(status, gin.H{
				"error":  c.Errors[0].Error(),
				"status": status,
			})

			c.Abort()
		}
	}
}
