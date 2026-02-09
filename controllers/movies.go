package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAppPage(c *gin.Context) {
	c.HTML(http.StatusOK, "app.html", BaseViewData("Movie Platform - Главная", "app"))
}

func GetMoviesPage(c *gin.Context) {
	c.HTML(http.StatusOK, "movies.html", BaseViewData("Movie Platform - Каталог", "movies"))
}

func GetMovieDetailsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "movie-details.html", BaseViewData("Movie Platform - Детали фильма", "movie"))
}

func GetProfilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", BaseViewData("Movie Platform - Профиль", "profile"))
}
