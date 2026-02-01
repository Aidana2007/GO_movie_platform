package routes

import (
	"github.com/Aidana2007/GO_movie_platform/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/movies", controllers.GetMovies())
	r.GET("/movies/:id", controllers.GetMovieById())
	r.POST("/movies", controllers.AddMovie())
	r.PUT("/movies/:id", controllers.UpdateMovie)
	r.DELETE("/movies/:id", controllers.DeleteMovie)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
}
