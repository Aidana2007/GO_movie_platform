package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Aidana2007/GO_movie_platform/internal/service"
	"github.com/Aidana2007/GO_movie_platform/pkg/utils"
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
	user := h.getCurrentUser(c)

	movies, err := h.movieService.GetAllMovies()
	if err != nil {
		log.Printf("Error fetching movies: %v", err)
		h.render(c, http.StatusInternalServerError, "error.html", gin.H{
			"User":  user,
			"Error": "Failed to fetch movies. Check movie records in DB (createdBy must be ObjectID).",
		})
		return
	}

	users, err := h.userService.SearchUsers("", 100) // Fetch all/many users
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		h.render(c, http.StatusInternalServerError, "error.html", gin.H{
			"User":  user,
			"Error": "Failed to fetch users",
		})
		return
	}

	h.render(c, http.StatusOK, "admin_panel.html", gin.H{
		"User":   user,
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

func (h *AdminHandler) render(c *gin.Context, status int, templateName string, data gin.H) {
	tmpl, err := template.New("base.html").ParseFiles(
		"../../../frontend/templates/base.html",
		"../../../frontend/templates/"+templateName,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template error: "+err.Error())
		return
	}

	c.Status(status)
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Template execution error: "+err.Error())
	}
}

func (h *AdminHandler) getCurrentUser(c *gin.Context) interface{} {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		return nil
	}

	user, err := h.userService.GetUserByID(userID.(primitive.ObjectID))
	if err != nil {
		return nil
	}

	return user
}
