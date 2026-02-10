package middleware

import (
	"strings"

	"github.com/Aidana2007/GO_movie_platform/backend/internal/repository"
	"github.com/Aidana2007/GO_movie_platform/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(jwtSecret string, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		cookie, err := c.Cookie("token")
		if err == nil {
			tokenString = cookie
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.Redirect(303, "/login")
				c.Abort()
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(401, gin.H{"error": "Invalid authorization header"})
				c.Abort()
				return
			}
			tokenString = parts[1]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.SetCookie("token", "", -1, "/", "", false, true)
			c.Redirect(303, "/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		role, err := userRepo.GetRoleByID(userID)
		if err != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set(utils.UserIDKey, userID)
		c.Set(utils.EmailKey, claims["email"].(string))
		c.Set(utils.RoleKey, role)

		c.Next()
	}
}
