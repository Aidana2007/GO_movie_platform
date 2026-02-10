package handler

import (
	"github.com/Aidana2007/GO_movie_platform/backend/internal/service"
	"github.com/Aidana2007/GO_movie_platform/backend/pkg/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService  *service.UserService
	movieService *service.MovieService
	authService  *service.AuthService
}

func NewUserHandler(userService *service.UserService, movieService *service.MovieService, authService *service.AuthService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		movieService: movieService,
		authService:  authService,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetWatchlist(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	movieIDs, err := h.userService.GetWatchlist(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get watchlist"})
		return
	}

	if len(movieIDs) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	movies, err := h.movieService.GetMoviesByIDs(movieIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *UserHandler) AddToWatchlist(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	movieID, err := primitive.ObjectIDFromHex(c.Param("movieId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	_, err = h.movieService.GetMovieByID(movieID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	if err := h.userService.AddToWatchlist(userID.(primitive.ObjectID), movieID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to watchlist"})
}

func (h *UserHandler) RemoveFromWatchlist(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	movieID, err := primitive.ObjectIDFromHex(c.Param("movieId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	if err := h.userService.RemoveFromWatchlist(userID.(primitive.ObjectID), movieID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from watchlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from watchlist"})
}
