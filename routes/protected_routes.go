package routes

import (
	"fmt"

	"github.com/Aidana2007/GO_movie_platform/controllers"
	"github.com/Aidana2007/GO_movie_platform/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(r *gin.Engine) {
	r.Use(middleware.AuthMiddleware)

	r.GET("/movies", controllers.GetMovies)
	r.GET("/movies/:id", controllers.GetMovieById)
	r.POST("/addmovie", controllers.AddMovie)
	r.PUT("/movies/:id", controllers.UpdateMovie)
	r.DELETE("/movies/:id", controllers.DeleteMovie)

	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
	}
}
