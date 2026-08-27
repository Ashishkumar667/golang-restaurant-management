package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Ashishkumar667/golang-restaurant-management/controllers"
)

func FoodRoutes(router *gin.Engine) {
	router.GET("/foods", controllers.GetAllFoods())
	router.GET("/foods/:food_id", controllers.GetFoodById())
	router.POST("/foods", controllers.CreateFood())
	router.PUT("/foods/:food_id", controllers.UpdateFood())
	router.DELETE("/foods/:food_id", controllers.DeleteFood())
}