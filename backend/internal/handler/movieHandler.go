package handler

import (
	"github.com/yerkebulan111/movie_smn/internal/model"
	"github.com/yerkebulan111/movie_smn/internal/service"
	"github.com/yerkebulan111/movie_smn/pkg/utils"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieHandler struct {
	movieService *service.MovieService
	authService  *service.AuthService
}

func NewMovieHandler(movieService *service.MovieService, authService *service.AuthService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
		authService:  authService,
	}
}

func (h *MovieHandler) GetMovies(c *gin.Context) {
	query := c.Query("search")
	genre := c.Query("genre")

	var movies []*model.Movie
	var err error

	if query != "" || genre != "" {
		sort := c.DefaultQuery("sort", "latest")
		minRatingStr := c.DefaultQuery("minRating", "0")
		minRating, _ := strconv.ParseFloat(minRatingStr, 64)
		movies, err = h.movieService.SearchMovies(query, genre, sort, minRating)
	} else {
		movies, err = h.movieService.GetAllMovies()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie, err := h.movieService.GetMovieByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req model.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Title == "" || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and description are required"})
		return
	}

	movie, err := h.movieService.CreateMovie(&req, userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var req model.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	movie, err := h.movieService.UpdateMovie(id, &req, userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	if err := h.movieService.DeleteMovie(id, userID.(primitive.ObjectID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
