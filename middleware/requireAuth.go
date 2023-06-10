package middleware

import (
	"belajar-golang-api-revisit-1/handler"
	"belajar-golang-api-revisit-1/user"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// get cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// decode and validate cookie
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		key := []byte("adfadfadf")

		return key, nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// find user with token
		var user user.User
		result := handler.Db.Where("email = ?", claims["sub"]).First(&user)
		if result.RowsAffected == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// attach to req
		c.Set("user", user)

		// response
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
