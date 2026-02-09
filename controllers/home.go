package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Aidana2007/GO_movie_platform/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetHome(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	movieCollection := database.GetCollection("movies")
	userCollection := database.GetCollection("users")
	reviewCollection := database.GetCollection("reviews")

	movieCount, _ := movieCollection.CountDocuments(ctx, bson.M{})
	userCount, _ := userCollection.CountDocuments(ctx, bson.M{})
	reviewCount, _ := reviewCollection.CountDocuments(ctx, bson.M{})

	data := BaseViewData("Movie Platform - Главная", "home")
	data["TotalMovies"] = movieCount
	data["TotalUsers"] = userCount
	data["TotalReviews"] = reviewCount

	c.HTML(http.StatusOK, "index.html", data)
}
