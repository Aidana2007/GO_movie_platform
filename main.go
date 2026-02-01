package main

import (
	"log"

	"github.com/Aidana2007/GO_movie_platform/controllers"
	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/Aidana2007/GO_movie_platform/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()
	controllers.InitWorker()
	controllers.StartMovieWorker()

	r := gin.Default()
	routes.RegisterRoutes(r)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
