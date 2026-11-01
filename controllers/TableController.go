package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Ashishkumar667/golang-restaurant-management/database"
	"github.com/Ashishkumar667/golang-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableControllerCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetAllTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cursor, err := tableControllerCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
			return
		}
		defer cursor.Close(ctx)

		var alltables []bson.M

		if err = cursor.All(ctx, &alltables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.JSON(http.StatusOK, alltables)
	}
}

func GetTableById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		tableId := c.Param("table_id")
		var table models.Table

		err := tableControllerCollection.FindOne(ctx, bson.M{"table_Id": tableId}).Decode(&table)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, table)
	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var table models.Table

		tableId := c.Param("table_id")
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if table.Number_of_guest != nil {
			updateObj = append(updateObj, bson.E{"number_of_guest", *table.Number_of_guest})
		}

		if table.Table_Number != nil {
			updateObj = append(updateObj, bson.E{"table_number", *table.Table_Number})
		}

		table.Updated_At = time.Now()
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"table_id": tableId}

		result, err := tableControllerCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "table was not updated"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(table)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.Created_At = time.Now()
		table.Updated_At = time.Now()
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		result, err := tableControllerCollection.InsertOne(ctx, table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "table was not created"})
			return
		}

		defer cancel()
		c.JSON(http.StatusCreated, result)
	}
}

func DeleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		tableId := c.Param("table_id")

		result, err := tableControllerCollection.DeleteOne(ctx, bson.M{"table_id": tableId})
		if err != nil {
			msg := fmt.Sprintf("table was not deleted")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
