package middleware

import (
	"net/http"

	"github.com/Aidana2007/GO_movie_platform/utils"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "no claims found"})
			c.Abort()
			return
		}

		userClaims := claims.(*utils.SignedDetails)

		if userClaims.Role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
