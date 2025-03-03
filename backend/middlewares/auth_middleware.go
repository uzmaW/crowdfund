package middlewares

import (
	"crowdfund/backend/services"
	"crowdfund/backend/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		tokenID, err := utils.ExtractTokenID(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if utils.IsTokenRevoked(tokenID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
			return
		}

		claims, err := utils.ValidateJWT(tokenString, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userService := services.NewUserService(c.MustGet("db").(*gorm.DB))
		user, err := userService.GetUserByID(uint(claims.UserID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
