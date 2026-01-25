package routes

import (
	controllers "magicstream/contollers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", controllers.Health)

	r.GET("/movies", controllers.GetMovies)
	r.POST("/movies", controllers.AddMovie)

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
}
