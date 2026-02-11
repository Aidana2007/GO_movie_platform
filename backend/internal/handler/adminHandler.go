package handler

import (
	"log"
	"net/http"

	"github.com/Aidana2007/GO_movie_platform/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminHandler struct {
	movieService  *service.MovieService
	userService   *service.UserService
	reviewService *service.ReviewService
}

func NewAdminHandler(movieService *service.MovieService, userService *service.UserService, reviewService *service.ReviewService) *AdminHandler {
	return &AdminHandler{
		movieService:  movieService,
		userService:   userService,
		reviewService: reviewService,
	}
}

func (h *AdminHandler) AdminDashboard(c *gin.Context) {
	movies, err := h.movieService.GetAllMovies()
	if err != nil {
		log.Printf("Error fetching movies: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to fetch movies"})
		return
	}

	users, err := h.userService.SearchUsers("", 100) // Fetch all/many users
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"Error": "Failed to fetch users"})
		return
	}

	c.HTML(http.StatusOK, "admin_panel.html", gin.H{
		"Movies": movies,
		"Users":  users,
	})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		log.Printf("Error deleting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User deleted successfully"})
}

func (h *AdminHandler) DeleteReviewAdmin(c *gin.Context) {
	reviewIDStr := c.Param("id")
	reviewID, err := primitive.ObjectIDFromHex(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.reviewService.DeleteReviewAdmin(reviewID); err != nil {
		log.Printf("Error deleting review (admin): %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
