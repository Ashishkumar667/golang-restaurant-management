package middlewares

import (
    "net/http"
    "strings"

    "github.com/Ashishkumar667/golang-restaurant-management/helpers"
    "github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
    return func(c *gin.Context) {

        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
            c.Abort() 
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be: Bearer <token>"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        claims, errMsg, err := helpers.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
            c.Abort()
            return
        }

        c.Set("email",      claims.Email)
        c.Set("first_name", claims.First_Name)
        c.Set("last_name",  claims.Last_Name)
        c.Set("user_id",    claims.User_id)

        c.Next()
    }
}