package main

import (
	"log"
	"os"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/Aidana2007/GO_movie_platform/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Load Go HTML templates
	r.LoadHTMLGlob("templates/*")

	// Static files
	r.Static("/static", "./static")

	// Register routes
	routes.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server: http://localhost:%s", port)

	r.Run(":" + port)
}
