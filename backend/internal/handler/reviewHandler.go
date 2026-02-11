package handler

import (
	"github.com/Aidana2007/GO_movie_platform/internal/model"
	"github.com/Aidana2007/GO_movie_platform/internal/service"
	"github.com/Aidana2007/GO_movie_platform/pkg/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
	authService   *service.AuthService
}

func NewReviewHandler(reviewService *service.ReviewService, authService *service.AuthService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		authService:   authService,
	}
}

func (h *ReviewHandler) GetMovieReviews(c *gin.Context) {
	movieID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	reviews, err := h.reviewService.GetMovieReviews(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": reviews})
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	movieID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.authService.GetUserByID(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	review, err := h.reviewService.CreateReview(
		movieID,
		&req,
		userID.(primitive.ObjectID),
		user.Username,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": review})
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	reviewID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := h.reviewService.DeleteReview(reviewID, userID.(primitive.ObjectID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Review deleted successfully"})
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID, exists := c.Get(utils.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	reviewID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	review, err := h.reviewService.UpdateReview(reviewID, &req, userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": review})
}
