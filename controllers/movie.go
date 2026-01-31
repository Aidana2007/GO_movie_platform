package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/Aidana2007/GO_movie_platform/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.GetCollection("movies")
var validate = validator.New()

// Movie processing worker
type MovieJob struct {
	Movie   models.Movie
	Action  string // "create", "log", "count"
}

var (
	movieJobQueue = make(chan MovieJob, 100)
	movieCounter  = struct {
		sync.RWMutex
		count int
	}{count: 0}
	workerStarted = false
	workerMutex   sync.Mutex
)

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies"})
		}

		c.JSON(http.StatusOK, movies)

	}
}

func GetMovieById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "there is no such movie"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

// StartBackgroundWorker starts the async movie processing worker
func StartBackgroundWorker() {
	workerMutex.Lock()
	defer workerMutex.Unlock()
	
	if workerStarted {
		return
	}
	workerStarted = true
	
	go func() {
		for job := range movieJobQueue {
			switch job.Action {
			case "create":
				// Async logging of movie creation
				logMovieCreation(job.Movie)
			case "count":
				// Increment thread-safe counter
				incrementMovieCounter()
			case "log":
				// General async logging
				logMovieEvent(job.Movie, "processed")
			}
		}
	}()
}

func logMovieCreation(movie models.Movie) {
	// Simulate async logging/processing
	time.Sleep(10 * time.Millisecond) // Simulate processing time
	// In production, this could write to a log file or external service
}

func logMovieEvent(movie models.Movie, event string) {
	// Async event logging
	time.Sleep(5 * time.Millisecond)
}

func incrementMovieCounter() {
	movieCounter.Lock()
	defer movieCounter.Unlock()
	movieCounter.count++
}

func GetMovieCount() int {
	movieCounter.RLock()
	defer movieCounter.RUnlock()
	return movieCounter.count
}

func AddMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var movie models.Movie

	// Bind JSON input
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input", "details": err.Error()})
		return
	}

	// Validate input
	if err := validate.Struct(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Check if movie with same imdb_id already exists
	var existingMovie models.Movie
	err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movie.ImdbID}).Decode(&existingMovie)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Movie with this imdb_id already exists"})
		return
	}
	// If error is not "no documents found", it's a real error
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing movie", "details": err.Error()})
		return
	}

	// Insert into MongoDB
	result, err := movieCollection.InsertOne(ctx, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie", "details": err.Error()})
		return
	}

	// Start background worker if not started
	StartBackgroundWorker()

	// Send job to background worker for async processing
	movieJobQueue <- MovieJob{
		Movie:  movie,
		Action: "create",
	}
	movieJobQueue <- MovieJob{
		Movie:  movie,
		Action: "count",
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Movie created successfully",
		"id":      result.InsertedID,
		"movie":   movie,
	})
}

// GetMovieStats returns statistics about movies (for testing async counter)
func GetMovieStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		StartBackgroundWorker() // Ensure worker is started
		c.JSON(http.StatusOK, gin.H{
			"total_processed": GetMovieCount(),
		})
	}
}
