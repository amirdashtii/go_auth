package middleware

import (
	"net/http"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	authService := service.NewAuthService(config)

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if ctx.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errors.ErrContextCancelled,
			})
			c.Abort()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.ErrMissingAuthHeader,
			})
			c.Abort()
			return
		}

		if len(token) > 7 && token[0:7] == "Bearer " {
			token = token[7:]
		}

		jwtSecret := config.JWT.Secret

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errors.ErrParseToken,
			})
			c.Abort()
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errors.ErrInvalidTokenClaims,
			})
			c.Abort()
			return
		}

		tokenType := claims["token_type"].(string)
		if tokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errors.ErrInvalidTokenType,
			})
			c.Abort()
			return
		}

		userID := claims["user_id"].(string)

		err = authService.ValidateToken(ctx, userID, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		roleClaim := claims["role"]
		var roleString string

		switch v := roleClaim.(type) {
		case float64:
			roleString = entities.RoleType(int(v)).String()
		case string:
			roleString = v
		default:
			roleString = "Unknown"
		}

		c.Set("user_id", userID)
		c.Set("role", roleString)
		c.Next()
	}
}
