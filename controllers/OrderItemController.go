package controllers

import (
	"context"
	"log"
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

type OrderItemPack struct {
	Table_id    *string
	order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetAllOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(ctx, bson.M{})

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItemById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orderId := c.Param("order_id")
		result, err := orderItemCollection.Find(
			ctx,
			bson.M{"order_id": orderId},
		)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date = time.Now()
		orderItemsToBeInserted := []interface{}{} //we can use similar to this []any{}{}

		order_id := OrderItemOrderCreator(order, ctx)

		for _, orderItem := range orderItemPack.order_items {
			orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_At = time.Now()
			orderItem.Updated_At = time.Now()
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		insertOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, insertOrderItems)
	}
}

// func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
// 	return func(c *gin.Context) {
// 		return c.JSON(http.StatusOK, OrderItems)
// 	}
// }
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.OrderItem

		orderItemId := c.Param("order_item_id")

		filter := bson.M{"order_item_id": orderItemId}

		var updateObject primitive.D

		if orderItem.Unit_price != nil {
			updateObject = append(updateObject, bson.E{Key: "unit_price", Value: *orderItem.Unit_price})
		}
		orderItem.Updated_At = time.Now()
		if orderItem.Quantity != nil {
			updateObject = append(updateObject, bson.E{Key: "quantity", Value: *orderItem.Quantity})
		}
		if orderItem.Food_id != nil {
			updateObject = append(updateObject, bson.E{Key: "food_id", Value: *orderItem.Food_id})
		}

		updateObject = append(updateObject, bson.E{Key: "updated_at", Value: orderItem.Updated_At})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObject},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order item update failed"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func DeleteOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
