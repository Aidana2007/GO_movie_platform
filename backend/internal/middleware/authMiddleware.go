package middleware

import (
	"github.com/yerkebulan111/movie_smn/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
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
			c.JSON(401, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		c.Set(utils.UserIDKey, userID)
		c.Set(utils.EmailKey, claims["email"].(string))

		c.Next()
	}
}
