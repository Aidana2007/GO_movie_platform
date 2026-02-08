package routes

import (
	"github.com/Aidana2007/GO_movie_platform/controllers"
	"github.com/Aidana2007/GO_movie_platform/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	public := r.Group("/")
	{
		public.POST("/register", controllers.RegisterUser)
		public.POST("/login", controllers.LoginUser)
		public.GET("/movies", controllers.GetMovies)
		public.GET("/movies/:id", controllers.GetMovieById)
		public.GET("/movies/:id/reviews", controllers.GetReviewsByMovie)
	}

	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.POST("/addmovie", controllers.AddMovie)
		admin.PUT("/movies/:id", controllers.UpdateMovie)
		admin.DELETE("/movies/:id", controllers.DeleteMovie)
	}

	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.POST("/movies/:id/reviews", controllers.AddReview)
		authenticated.DELETE("/reviews/:id", controllers.DeleteReview)
	}
}
