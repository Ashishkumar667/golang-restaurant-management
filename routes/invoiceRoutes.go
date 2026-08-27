package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/Ashishkumar667/golang-restaurant-management/controllers"
)

func InvoiceRoutes(router *gin.Engine) {
	router.GET("/invoices", controllers.GetInvoices())
	router.GET("/invoices/:invoice_id", controllers.GetInvoiceById())
	router.POST("/invoices", controllers.CreateInvoice())
	router.PUT("/invoices/:invoice_id", controllers.UpdateInvoice())
	router.DELETE("/invoices/:invoice_id", controllers.DeleteInvoice())
}