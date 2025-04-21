package middleware

import (
	"net/http"
	"product-management/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "authorization header is required",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Check format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "invalid authorization header format",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Load config for secret key
		cfg, err := config.LoadConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "failed to load configuration",
				"status": http.StatusInternalServerError,
			})
			c.Abort()
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure token uses HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "invalid or expired token",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "invalid token claims",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Extract specific claims: user_id, email, role
		userIDFloat, okID := claims["user_id"].(float64)
		email, okEmail := claims["email"].(string)
		role, okRole := claims["role"].(string)

		if !okID || !okEmail || !okRole {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "missing or invalid claim fields",
				"status": http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Set into context
		c.Set("userID", uint(userIDFloat))
		c.Set("email", email)
		c.Set("role", role)

		c.Next()
	}
}
