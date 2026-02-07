package middleware

import (
	"net/http"

	"github.com/Aidana2007/GO_movie_platform/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token, err := utils.GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
		c.Abort()
		return
	}
	claims, err := utils.ValidateToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.Set("userId", claims.UserId)
	c.Set("role", claims.Role)

	c.Next()
}
