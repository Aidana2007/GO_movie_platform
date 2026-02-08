package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/Aidana2007/GO_movie_platform/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var validate = validator.New()

func getMovieCollection() *mongo.Collection {
	return database.GetCollection("movies")
}

func GetMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var movies []models.Movie

	cursor, err := getMovieCollection().Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movies"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func GetMovieById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var movie models.Movie
	err = getMovieCollection().FindOne(ctx, bson.M{"_id": objID}).Decode(&movie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func AddMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie data"})
		return
	}

	if err := validate.Struct(movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genreCollection := database.GetCollection("genres")
	for _, gid := range movie.GenreIDs {
		count, _ := genreCollection.CountDocuments(ctx, bson.M{"genre_id": gid})
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre_id: " + strconv.Itoa(gid)})
			return
		}
	}

	result, err := getMovieCollection().InsertOne(ctx, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add movie"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "movie added successfully",
		"id":      result.InsertedID,
	})
}

func UpdateMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var updated models.Movie
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":       updated.Title,
			"description": updated.Description,
			"genre_ids":   updated.GenreIDs,
			"year":        updated.Year,
			"poster_path": updated.PosterPath,
			"video_path":  updated.VideoPath,
		},
	}

	result, err := getMovieCollection().UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie updated successfully"})
}

func DeleteMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	reviewCollection := database.GetCollection("reviews")
	_, err = reviewCollection.DeleteMany(ctx, bson.M{"movie_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete related reviews"})
		return
	}

	result, err := getMovieCollection().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie and its reviews deleted successfully"})
}
