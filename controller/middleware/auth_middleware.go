package middleware

import (
	"net/http"
	"strings"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	authService := service.NewAuthService()

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := authHeader
		if len(authHeader) > 7 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
			tokenString = authHeader[7:]
		}

		user, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid or expired token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		userValue, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		user, ok := userValue.(*entities.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve user information",
			})
			c.Abort()
			return
		}

		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin privileges required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
