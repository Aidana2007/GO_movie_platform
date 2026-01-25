package main

import (
	"log"
	"magicstream/database"
	"magicstream/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "SEMovie running",
		})
	})

	routes.RegisterRoutes(r)
	database.ConnectMongo()

	log.Println("Server running on :8080")
	r.Run(":8080")
}
