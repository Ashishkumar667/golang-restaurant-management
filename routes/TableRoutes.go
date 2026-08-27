package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Ashishkumar667/golang-restaurant-management/controllers"
)

func TableRoutes(router *gin.Engine) {
	router.GET("/tables", controllers.GetAllTables())
	router.GET("/tables/:table_id", controllers.GetTableById())
	router.POST("/tables", controllers.CreateTable())
	router.PUT("/tables/:table_id", controllers.UpdateTable())
	router.DELETE("/tables/:table_id", controllers.DeleteTable())
}