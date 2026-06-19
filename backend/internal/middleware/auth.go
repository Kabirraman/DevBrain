package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})

			c.Abort()
			return
		}

		tokenString := authHeader[7:] // Remove "Bearer "

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			return []byte(os.Getenv("JWT_SECRET")), nil

		})

		if err != nil || !token.Valid {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})

			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set(
			"userID",
			claims["user_id"],
		)

		c.Next()
	}
}