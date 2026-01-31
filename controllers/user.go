package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Register endpoint works",
	})
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Login endpoint works",
	})
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "SEMovie API running",
	})
}
