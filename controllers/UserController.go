package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Ashishkumar667/golang-restaurant-management/database"
	"github.com/Ashishkumar667/golang-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "User")

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := UserCollection.Find(ctx, bson.M{})
		defer result.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch Users data"})
			return
		}

		var user []models.User

		if err := result.All(ctx, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch users data"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		UserID := c.Param("user_id")
		var user models.User
		err := UserCollection.FindOne(ctx, bson.M{"user_id": UserID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, user)

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}

		hashedPassword, err := HashPassword(*user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
			return
		}
		user.Password = &hashedPassword

		user.Created_At = time.Now()
		user.Updated_At = time.Now()
		user.ID = primitive.NewObjectID()

		result, err := UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var loginData models.User
		var foundUser models.User

		if err := c.BindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		isValid, msg := VerifyPassword(*loginData.Password, *foundUser.Password)
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login successful",
			"user_id": foundUser.ID,
			"email":   foundUser.Email,
		})
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(userPassword string, ProvidedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(ProvidedPassword), []byte(userPassword))
	if err != nil {
		return false, "password is incorrect"
	}

	return true, ""
}
