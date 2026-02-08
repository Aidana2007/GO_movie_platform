package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/Aidana2007/GO_movie_platform/models"
	"github.com/Aidana2007/GO_movie_platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func getReviewCollection() *mongo.Collection {
	return database.GetCollection("reviews")
}

func AddReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userClaims := claims.(*utils.SignedDetails)

	movieIDStr := c.Param("id")
	movieID, err := bson.ObjectIDFromHex(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	movieCollection := getMovieCollection()
	count, err := movieCollection.CountDocuments(ctx, bson.M{"_id": movieID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check movie"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.MovieID = movieID
	review.UserID = userClaims.UserId
	review.CreatedAt = time.Now().Format(time.RFC3339)

	count, err = getReviewCollection().CountDocuments(ctx, bson.M{
		"movie_id": movieID,
		"user_id":  userClaims.UserId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing review"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "review already exists for this movie"})
		return
	}

	result, err := getReviewCollection().InsertOne(ctx, review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "review added successfully",
		"id":      result.InsertedID,
	})
}

func GetReviewsByMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	movieIDStr := c.Param("id")
	movieID, err := bson.ObjectIDFromHex(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var reviews []models.Review
	cursor, err := getReviewCollection().Find(ctx, bson.M{"movie_id": movieID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reviews"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &reviews); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reviews,
	})
}

func DeleteReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userClaims := claims.(*utils.SignedDetails)

	reviewIDStr := c.Param("id")
	reviewID, err := bson.ObjectIDFromHex(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	var review models.Review
	err = getReviewCollection().FindOne(ctx, bson.M{"_id": reviewID}).Decode(&review)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	if review.UserID != userClaims.UserId {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own review"})
		return
	}

	result, err := getReviewCollection().DeleteOne(ctx, bson.M{"_id": reviewID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}
