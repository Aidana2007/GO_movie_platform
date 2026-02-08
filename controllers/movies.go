package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAppPage(c *gin.Context) {
	c.HTML(http.StatusOK, "app.html", nil)
}

func GetMoviesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "movies.html", nil)
}

func GetMovieDetailsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "movie-details.html", nil)
}

func GetProfilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", nil)
}
