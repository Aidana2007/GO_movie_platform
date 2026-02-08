package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type HomeData struct {
	TotalMovies  int64
	TotalUsers   int64
	TotalReviews int64
}

func GetHome(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var data HomeData

	movieCollection := database.GetCollection("movies")
	userCollection := database.GetCollection("users")
	reviewCollection := database.GetCollection("reviews")

	movieCount, _ := movieCollection.CountDocuments(ctx, bson.M{})
	userCount, _ := userCollection.CountDocuments(ctx, bson.M{})
	reviewCount, _ := reviewCollection.CountDocuments(ctx, bson.M{})

	data.TotalMovies = movieCount
	data.TotalUsers = userCount
	data.TotalReviews = reviewCount

	c.HTML(http.StatusOK, "index.html", data)
}
