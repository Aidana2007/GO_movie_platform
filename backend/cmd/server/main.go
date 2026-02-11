package main

import (
	"html/template"
	"log"
	"net/url"
	"strings"

	"github.com/Aidana2007/GO_movie_platform/internal/config"
	"github.com/Aidana2007/GO_movie_platform/internal/handler"
	"github.com/Aidana2007/GO_movie_platform/internal/middleware"
	"github.com/Aidana2007/GO_movie_platform/internal/repository"
	"github.com/Aidana2007/GO_movie_platform/internal/service"
	"github.com/gin-gonic/gin"
)

func extractYouTubeEmbed(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	if strings.Contains(u.Host, "youtu.be") {
		videoID := strings.TrimPrefix(u.Path, "/")
		return "https://www.youtube.com/embed/" + videoID
	}

	if strings.Contains(u.Host, "youtube.com") {
		videoID := u.Query().Get("v")
		if videoID != "" {
			return "https://www.youtube.com/embed/" + videoID
		}
	}

	return urlStr
}

func main() {
	cfg := config.LoadConfig()

	client, err := config.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer config.DisconnectMongoDB(client)

	db := client.Database(cfg.DBName)

	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	friendRepo := repository.NewFriendRepository(db)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	movieService := service.NewMovieService(movieRepo)
	reviewService := service.NewReviewService(reviewRepo, movieRepo)
	friendService := service.NewFriendService(friendRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	movieHandler := handler.NewMovieHandler(movieService, authService)
	reviewHandler := handler.NewReviewHandler(reviewService, authService)
	userHandler := handler.NewUserHandler(userService, movieService, authService)
	pageHandler := handler.NewPageHandler(movieService, authService, userService, friendService)
	friendHandler := handler.NewFriendHandler(friendService, userService)
	adminHandler := handler.NewAdminHandler(movieService, userService, reviewService)

	auth := middleware.AuthMiddleware(cfg.JWTSecret, userRepo)

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"extractYouTubeEmbed": extractYouTubeEmbed,
	})
	r.LoadHTMLGlob("../../../frontend/templates/*")
	r.Static("/static", "../../../frontend/static")

	r.GET("/", pageHandler.HomePage)
	r.GET("/movies", pageHandler.MoviesPage)
	r.GET("/movie/:id", pageHandler.MovieDetailsPage)
	r.GET("/login", pageHandler.LoginPage)
	r.GET("/register", pageHandler.RegisterPage)

	authPages := r.Group("/")
	authPages.Use(auth)
	{
		authPages.GET("/watchlist", pageHandler.WatchlistPage)
		authPages.GET("/profile", pageHandler.ProfilePage)
		authPages.GET("/users", pageHandler.UsersPage)
		authPages.GET("/user/:id", pageHandler.UserProfilePage)

		adminPages := authPages.Group("/")
		adminPages.Use(middleware.RequireAdmin())
		{
			adminPages.GET("/admin", adminHandler.AdminDashboard)
		}
	}

	api := r.Group("/api")
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)

	api.GET("/movies", movieHandler.GetMovies)
	api.GET("/movies/:id", movieHandler.GetMovieByID)
	api.GET("/movies/:id/reviews", reviewHandler.GetMovieReviews)

	apiAuth := api.Group("/")
	apiAuth.Use(auth)
	{

		adminMovies := apiAuth.Group("/")
		adminMovies.Use(middleware.RequireAdmin())
		{
			adminMovies.POST("/movies", movieHandler.CreateMovie)
			adminMovies.PUT("/movies/:id", movieHandler.UpdateMovie)
			adminMovies.DELETE("/movies/:id", movieHandler.DeleteMovie)
			adminMovies.DELETE("/admin/users/:id", adminHandler.DeleteUser)
			adminMovies.DELETE("/admin/reviews/:id", adminHandler.DeleteReviewAdmin)
		}

		apiAuth.POST("/movies/:id/reviews", reviewHandler.CreateReview)
		apiAuth.PUT("/reviews/:id", reviewHandler.UpdateReview)
		apiAuth.DELETE("/reviews/:id", reviewHandler.DeleteReview)

		modReviews := apiAuth.Group("/")
		modReviews.Use(middleware.RequireModeratorOrAdmin())
		{
			// Moderator specific routes if any
		}

		apiAuth.GET("/user/watchlist", userHandler.GetWatchlist)
		apiAuth.POST("/user/watchlist/:movieId", userHandler.AddToWatchlist)
		apiAuth.DELETE("/user/watchlist/:movieId", userHandler.RemoveFromWatchlist)
		apiAuth.GET("/user/profile", userHandler.GetProfile)
	}

	friends := api.Group("/user")
	friends.Use(auth)
	{
		friends.POST("/friend-request", friendHandler.SendFriendRequest)
		friends.GET("/friend-requests", friendHandler.GetFriendRequests)
		friends.GET("/sent-requests", friendHandler.GetSentRequests)
		friends.POST("/friend-request/:id/accept", friendHandler.AcceptFriendRequest)
		friends.POST("/friend-request/:id/reject", friendHandler.RejectFriendRequest)
		friends.DELETE("/friend-request/:id", friendHandler.CancelFriendRequest)
		friends.GET("/friends", friendHandler.GetFriends)
		friends.DELETE("/friends/:id", friendHandler.RemoveFriend)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
