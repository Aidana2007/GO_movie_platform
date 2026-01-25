package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GetMovies works"})
}

func AddMovie(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "AddMovie works"})
}
