package routes

import (
	"github.com/Aidana2007/GO_movie_platform/controllers"
	"github.com/Aidana2007/GO_movie_platform/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(r *gin.Engine) {
	r.Use(middleware.AuthMiddleware)

	r.GET("/movies/:id", controllers.GetMovieById)
	r.POST("/addmovie", controllers.AddMovie)

}
