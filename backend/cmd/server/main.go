package main

import (
	"github.com/yerkebulan111/movie_smn/internal/config"
	"github.com/yerkebulan111/movie_smn/internal/handler"
	"github.com/yerkebulan111/movie_smn/internal/middleware"
	"github.com/yerkebulan111/movie_smn/internal/repository"
	"github.com/yerkebulan111/movie_smn/internal/service"
	"html/template"
	"log"

	"net/url"
	"strings"

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
	authPages.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		authPages.GET("/watchlist", pageHandler.WatchlistPage)
		authPages.GET("/profile", pageHandler.ProfilePage)
	}

	api := r.Group("/api")

	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)

	api.GET("/movies", movieHandler.GetMovies)
	api.GET("/movies/:id", movieHandler.GetMovieByID)

	apiAuth := api.Group("/")
	apiAuth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		apiAuth.POST("/movies", movieHandler.CreateMovie)
		apiAuth.PUT("/movies/:id", movieHandler.UpdateMovie)
		apiAuth.DELETE("/movies/:id", movieHandler.DeleteMovie)

		apiAuth.POST("/movies/:id/reviews", reviewHandler.CreateReview)
		apiAuth.DELETE("/reviews/:id", reviewHandler.DeleteReview)

		apiAuth.GET("/user/watchlist", userHandler.GetWatchlist)
		apiAuth.POST("/user/watchlist/:movieId", userHandler.AddToWatchlist)
		apiAuth.DELETE("/user/watchlist/:movieId", userHandler.RemoveFromWatchlist)
		apiAuth.GET("/user/profile", userHandler.GetProfile)
	}

	friends := api.Group("/user")
	friends.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		friends.POST("/friend-request", friendHandler.SendFriendRequest)
		friends.GET("/friend-requests", friendHandler.GetFriendRequests)
		friends.POST("/friend-request/:id/accept", friendHandler.AcceptFriendRequest)
		friends.POST("/friend-request/:id/reject", friendHandler.RejectFriendRequest)
		friends.GET("/friends", friendHandler.GetFriends)
		friends.DELETE("/friends/:id", friendHandler.RemoveFriend)
	}

	api.GET("/users", middleware.AuthMiddleware(cfg.JWTSecret), func(c *gin.Context) {
		pageHandler.UsersPage(c)
	})
	api.GET("/user/:id", pageHandler.UserProfilePage)

	api.GET("/movies/:id/reviews", reviewHandler.GetMovieReviews)

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
