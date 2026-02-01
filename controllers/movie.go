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
)

var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieCollection := database.GetCollection("movies")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movies := make([]models.Movie, 0)

		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies"})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovieById() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieCollection := database.GetCollection("movies")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		idParam := c.Param("id")

		objID, err := bson.ObjectIDFromHex(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie id"})
			return
		}

		var movie models.Movie
		err = movieCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieCollection := database.GetCollection("movies")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie data"})
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
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre_id: " + strconv.Itoa(gid)})
				return
			}
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		MovieWorkerChan <- movie

		c.JSON(http.StatusCreated, result)
	}
}

func UpdateMovie(c *gin.Context) {
	movieCollection := database.GetCollection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	objID, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid movie id"})
		return
	}

	var updatedMovie models.Movie
	if err := c.ShouldBindJSON(&updatedMovie); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updateData, _ := bson.Marshal(updatedMovie)
	var updateMap bson.M
	bson.Unmarshal(updateData, &updateMap)

	delete(updateMap, "_id")

	update := bson.M{
		"$set": updateMap,
	}

	result, err := movieCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(404, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Movie updated successfully"})
}

func DeleteMovie(c *gin.Context) {
	movieCollection := database.GetCollection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")

	objID, err := bson.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid movie id"})
		return
	}

	reviewCollection := database.GetCollection("reviews")
	_, err = reviewCollection.DeleteMany(ctx, bson.M{"movie_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete related reviews"})
		return
	}

	result, err := movieCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(404, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Movie and its reviews deleted successfully"})
}
