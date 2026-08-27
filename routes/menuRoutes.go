package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Ashishkumar667/golang-restaurant-management/controllers"
)

func MenuRoutes(router *gin.Engine) {
	router.GET("/menus", controllers.GetAllMenus())
	router.GET("/menus/:menu_id", controllers.GetMenuById())
	router.POST("/menus", controllers.CreateMenu())
	router.PUT("/menus/:menu_id", controllers.UpdateMenu())
	router.DELETE("/menus/:menu_id", controllers.DeleteMenu())
}