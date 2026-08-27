package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Ashishkumar667/golang-restaurant-management/controllers"
)

func OrderRoutes(router *gin.Engine) {
	router.GET("/orders", controllers.GetAllorders())
	router.GET("/orders/:order_id", controllers.GetOrderById())
	router.POST("/orders", controllers.CreateOrder())
	router.POST("/orders/:order_id", controllers.UpdateOrder())
}