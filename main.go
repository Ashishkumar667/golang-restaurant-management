package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/Ashishkumar667/golang-restaurant-management/database"
	"github.com/Ashishkumar667/golang-restaurant-management/routes"
	// "github.com/Ashishkumar667/golang-restaurant-management/controllers"
	"github.com/Ashishkumar667/golang-restaurant-management/middlewares"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middlewares.Authentication())


	routes.FoodRoutes(router)
	routes.OrderRoutes(router)
	routes.PaymentRoutes(router)
	routes.TableRoutes(router)
	routes.MenuRoutes(router)
    routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)

}