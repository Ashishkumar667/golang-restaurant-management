package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Ashishkumar667/golang-restaurant-management/controllers"
)

func PaymentRoutes(router *gin.Engine) {
	router.GET("/payments", controllers.GetAllPayments())
	router.GET("/payments/:payment_id", controllers.GetPaymentById())
	router.POST("/payments", controllers.CreatePayment())
	router.PUT("/payments/:payment_id", controllers.UpdatePayment())
}