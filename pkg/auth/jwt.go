package auth

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GetUserID(c *gin.Context) int {
	tokenString := c.GetHeader("Authorization") // Bearer <token_string>
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return int(claims["user_id"].(float64))
	}
	return 0
}
