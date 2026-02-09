package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
)

const appName = "MOVIE"

func BaseViewData(title, page string) gin.H {
	return gin.H{
		"Title":   title,
		"AppName": appName,
		"Year":    time.Now().Year(),
		"Page":    page,
	}
}
