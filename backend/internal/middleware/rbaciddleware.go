package middleware

import (
	"net/http"

	"github.com/Aidana2007/GO_movie_platform/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(utils.RoleKey)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role information unavailable"})
			c.Abort()
			return
		}

		if role.(string) != utils.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireModeratorOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(utils.RoleKey)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role information unavailable"})
			c.Abort()
			return
		}

		roleStr := role.(string)
		if roleStr != utils.RoleModerator && roleStr != utils.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Moderator or admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
