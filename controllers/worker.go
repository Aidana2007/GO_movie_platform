package controllers

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/Aidana2007/GO_movie_platform/models"
)

var MovieWorkerChan chan models.Movie

var processedCount int64

func StartMovieWorker() {
	logFile, err := os.OpenFile("background_jobs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	go func() {
		logger := log.New(logFile, "MOVIE_WORKER: ", log.LstdFlags)
		logger.Println("Movie worker started")
		for movie := range MovieWorkerChan {
			atomic.AddInt64(&processedCount, 1)
			count := atomic.LoadInt64(&processedCount)
			msg := fmt.Sprintf("Processing movie: %s (%s). Total processed: %d", movie.Title, movie.ImdbID, count)
			logger.Println(msg)
			fmt.Println(msg)
			processMovie(movie)
		}
	}()
}

func processMovie(movie models.Movie) {
	time.Sleep(2 * time.Second)
	log.Printf("Processed movie: %s", movie.Title)
}
