package routes

import (
	"fmt"

	"github.com/Aidana2007/GO_movie_platform/controllers"
	"github.com/Aidana2007/GO_movie_platform/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUnprotectedRoutes(r *gin.Engine) {
	r.Use(middleware.AuthMiddleware)

	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
	}
}
