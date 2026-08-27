package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Ashishkumar667/golang-restaurant-management/controllers"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/users", controllers.GetAllUsers())
	router.GET("/users/:user_id", controllers.GetUserById())
	router.POST("/users/signup", controllers.SignUp())
	router.POST("/users/login", controllers.Login())
}