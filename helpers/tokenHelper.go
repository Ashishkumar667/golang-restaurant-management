package helpers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Ashishkumar667/golang-restaurant-management/database"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

type Claims struct {
	Email      string `json:"email"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	User_id    string `json:"user_id"`
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateTokens(email, firstName, lastName, userID string) (string, string, error) {

	claims := &Claims{
		Email:      email,
		First_Name: firstName,
		Last_Name:  lastName,
		User_id:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (*Claims, string, error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, "token is invalid", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, "token is invalid", fmt.Errorf("invalid claims")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, "token is expired", fmt.Errorf("token expired")
	}

	return claims, "", nil
}

func UpdateTokens(token, refreshToken, userID string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"refresh_token": refreshToken,
				"updated_at":    time.Now(),
			},
		},
		options.Update().SetUpsert(true),
	)

	return err
}
