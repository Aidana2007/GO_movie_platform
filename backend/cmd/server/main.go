package main

import (
	"html/template"
	"log"
	"net/url"
	"strings"

	"github.com/Aidana2007/GO_movie_platform/backend/internal/config"
	"github.com/Aidana2007/GO_movie_platform/backend/internal/handler"
	"github.com/Aidana2007/GO_movie_platform/backend/internal/middleware"
	"github.com/Aidana2007/GO_movie_platform/backend/internal/repository"
	"github.com/Aidana2007/GO_movie_platform/backend/internal/service"
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
	// Load configuration.
	cfg := config.LoadConfig()

	// MongoDB connection.
	client, err := config.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer config.DisconnectMongoDB(client)

	db := client.Database(cfg.DBName)

	// Repositories.
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	friendRepo := repository.NewFriendRepository(db)

	// Services.
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	movieService := service.NewMovieService(movieRepo)
	reviewService := service.NewReviewService(reviewRepo, movieRepo)
	friendService := service.NewFriendService(friendRepo, userRepo)

	// Handlers.
	authHandler := handler.NewAuthHandler(authService)
	movieHandler := handler.NewMovieHandler(movieService, authService)
	reviewHandler := handler.NewReviewHandler(reviewService, authService)
	userHandler := handler.NewUserHandler(userService, movieService, authService)
	pageHandler := handler.NewPageHandler(movieService, authService, userService, friendService)
	friendHandler := handler.NewFriendHandler(friendService, userService)

	// Convenience: build the auth middleware once, bound to userRepo.
	// AuthMiddleware now also injects the user's role from MongoDB.
	auth := middleware.AuthMiddleware(cfg.JWTSecret, userRepo)

	// Gin router.
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

	// Page routes that require a valid session.
	authPages := r.Group("/")
	authPages.Use(auth)
	{
		authPages.GET("/watchlist", pageHandler.WatchlistPage)
		authPages.GET("/profile", pageHandler.ProfilePage)
	}

	// ----------------------------------------------------------------
	// API routes
	// ----------------------------------------------------------------
	api := r.Group("/api")

	// Public auth endpoints.
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)

	// Public read-only movie/review endpoints.
	api.GET("/movies", movieHandler.GetMovies)
	api.GET("/movies/:id", movieHandler.GetMovieByID)
	api.GET("/movies/:id/reviews", reviewHandler.GetMovieReviews)
	api.GET("/user/:id", pageHandler.UserProfilePage)

	// ----------------------------------------------------------------
	// Authenticated API group
	// Requires: valid JWT + role fetched from MongoDB
	// ----------------------------------------------------------------
	apiAuth := api.Group("/")
	apiAuth.Use(auth)
	{
		// ----------------------------------------------------------
		// Movie CRUD — admin only
		// POST   /api/movies
		// PUT    /api/movies/:id
		// DELETE /api/movies/:id
		// ----------------------------------------------------------
		adminMovies := apiAuth.Group("/")
		adminMovies.Use(middleware.RequireAdmin())
		{
			adminMovies.POST("/movies", movieHandler.CreateMovie)
			adminMovies.PUT("/movies/:id", movieHandler.UpdateMovie)
			adminMovies.DELETE("/movies/:id", movieHandler.DeleteMovie)
		}

		// ----------------------------------------------------------
		// Review management
		// POST   /api/movies/:id/reviews  — any authenticated user
		// DELETE /api/reviews/:id         — moderator or admin only
		// ----------------------------------------------------------
		apiAuth.POST("/movies/:id/reviews", reviewHandler.CreateReview)

		modReviews := apiAuth.Group("/")
		modReviews.Use(middleware.RequireModeratorOrAdmin())
		{
			modReviews.DELETE("/reviews/:id", reviewHandler.DeleteReview)
		}

		// ----------------------------------------------------------
		// User self-service (watchlist, profile)
		// ----------------------------------------------------------
		apiAuth.GET("/user/watchlist", userHandler.GetWatchlist)
		apiAuth.POST("/user/watchlist/:movieId", userHandler.AddToWatchlist)
		apiAuth.DELETE("/user/watchlist/:movieId", userHandler.RemoveFromWatchlist)
		apiAuth.GET("/user/profile", userHandler.GetProfile)
	}

	// ----------------------------------------------------------------
	// Friend system — authenticated users only
	// ----------------------------------------------------------------
	friends := api.Group("/user")
	friends.Use(auth)
	{
		friends.POST("/friend-request", friendHandler.SendFriendRequest)
		friends.GET("/friend-requests", friendHandler.GetFriendRequests)
		friends.POST("/friend-request/:id/accept", friendHandler.AcceptFriendRequest)
		friends.POST("/friend-request/:id/reject", friendHandler.RejectFriendRequest)
		friends.GET("/friends", friendHandler.GetFriends)
		friends.DELETE("/friends/:id", friendHandler.RemoveFriend)
	}

	// Users listing page (authenticated).
	api.GET("/users", auth, func(c *gin.Context) {
		pageHandler.UsersPage(c)
	})

	// Start server.
	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
