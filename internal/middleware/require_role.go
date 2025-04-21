package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Access denied: missing role",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Access denied: invalid role format",
				"status": http.StatusForbidden,
			})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if strings.EqualFold(userRole, allowed) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":  "Access denied: insufficient permissions",
			"status": http.StatusForbidden,
		})
		c.Abort()
	}
}
