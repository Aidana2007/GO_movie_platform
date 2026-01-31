package routes

import (
	controllers "github.com/Aidana2007/GO_movie_platform/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", controllers.Health)

	r.GET("/movies", controllers.GetMovies())
	r.GET("/movies/:imdb_id", controllers.GetMovieById())
	r.GET("/movies/stats", controllers.GetMovieStats())

	r.POST("/movies", controllers.AddMovie())

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
}
