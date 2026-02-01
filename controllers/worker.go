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

func InitWorker() {
	MovieWorkerChan = make(chan models.Movie)
}

func StartMovieWorker() {
	logFile, _ := os.OpenFile("background_jobs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	go func() {
		logger := log.New(logFile, "MOVIE_WORKER: ", log.LstdFlags)
		logger.Println("Movie worker started")

		for movie := range MovieWorkerChan {
			atomic.AddInt64(&processedCount, 1)
			count := atomic.LoadInt64(&processedCount)

			msg := fmt.Sprintf("Processing movie: %s | Total processed: %d", movie.Title, count)
			logger.Println(msg)
			fmt.Println(msg)

			time.Sleep(2 * time.Second)
		}
	}()
}
