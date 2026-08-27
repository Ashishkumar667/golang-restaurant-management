package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Ashishkumar667/golang-restaurant-management/controllers"
)

func OrderItemRoutes(router *gin.Engine) {
	router.GET("/order-items", controllers.GetAllOrderItems())
	router.GET("/order-items/:order_item_id", controllers.GetOrderItemById())
	router.POST("/order-items", controllers.CreateOrderItem())
	router.PUT("/order-items/:order_item_id", controllers.UpdateOrderItem())
	router.DELETE("/order-items/:order_item_id", controllers.DeleteOrderItem())
}