package main

import (
	"log"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/Aidana2007/GO_movie_platform/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.RegisterRoutes(r)
	database.ConnectDB()

	log.Println("Server running on :8080")
	r.Run(":8080")
}
