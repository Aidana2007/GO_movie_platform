package routes

import (
	"github.com/Aidana2007/GO_movie_platform/controllers"
	"github.com/Aidana2007/GO_movie_platform/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// Go Template routes (requirement)
	public := r.Group("/")
	{
		public.GET("/", controllers.GetHome)
		public.GET("/app", controllers.GetAppPage)
		public.GET("/movies", controllers.GetMoviesPage)
		public.GET("/movie/:id", controllers.GetMovieDetailsPage)
		public.GET("/profile", controllers.GetProfilePage)
	}

	// API routes
	api := r.Group("/api")
	{
		// Auth
		api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.LoginUser)

		// Movies (public)
		api.GET("/movies", controllers.GetMovies)
		api.GET("/movies/:id", controllers.GetMovieById)
		api.GET("/movies/:id/reviews", controllers.GetReviewsByMovie)

		// Admin routes
		admin := api.Group("/")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			admin.POST("/movies", controllers.AddMovie)
			admin.PUT("/movies/:id", controllers.UpdateMovie)
			admin.DELETE("/movies/:id", controllers.DeleteMovie)
		}

		// Authenticated routes
		authenticated := api.Group("/")
		authenticated.Use(middleware.AuthMiddleware())
		{
			authenticated.POST("/movies/:id/reviews", controllers.AddReview)
			authenticated.DELETE("/reviews/:id", controllers.DeleteReview)
		}
	}
}
