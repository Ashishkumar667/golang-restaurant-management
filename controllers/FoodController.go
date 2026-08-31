package controllers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/Ashishkumar667/golang-restaurant-management/database"
	"github.com/Ashishkumar667/golang-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

func GetAllFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := foodCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching foods"})
			return
		}
		defer cursor.Close(ctx)

		var foods []models.Food

		if err = cursor.All(ctx, &foods); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching foods"})
			return
		}

		c.JSON(http.StatusOK, foods)
	}
}

func GetFoodById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching the food item"})
			return
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("Menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		food.Created_At = time.Now()
		food.Updated_At = time.Now()
		food.ID = primitive.NewObjectID()
		FOODID := food.ID.Hex()
		food.Food_id = &FOODID
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		foodId := c.Param("food_id")

		foodObjectID, err := primitive.ObjectIDFromHex(foodId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid food id",
			})
			return
		}

		var food models.Food

		if err := c.ShouldBindJSON(&food); err != nil { //read json payload from client
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		food.Updated_At = time.Now()

		var updatedFood models.Food

		updateErr := foodCollection.FindOneAndUpdate(
			ctx,
			bson.M{
				"_id": foodObjectID,
			},

			bson.M{
				"$set": bson.M{
					"name":       food.Name,
					"price":      food.Price,
					"food_image": food.Food_image,
					"updated_at": food.Updated_At,
				},
			},
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&updatedFood)

		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update food"})
			return
		}

		c.JSON(http.StatusOK, updatedFood)
	}
}

func DeleteFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		foodId := c.Param("food_id")

		foodObjectID, err := primitive.ObjectIDFromHex(foodId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid food id",
			})
			return
		}

		result, err := foodCollection.DeleteOne(ctx, bson.M{"food_id": foodObjectID})

		if err != nil {
			msg := fmt.Sprintf("Could not delete tbis food")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "This food does not exist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "food Deleted Successfully"})
	}
}
