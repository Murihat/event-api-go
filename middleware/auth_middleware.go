package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequiredAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No Access, token is not found!",
			})

			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context.Set("userID", int(claims["sub"].(float64)))
			context.Next()
		} else {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
		}
	}
}
