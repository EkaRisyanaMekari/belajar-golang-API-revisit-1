package handler

import (
	"belajar-golang-api-revisit-1/user"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Read request body
	type Body struct {
		Email    string
		Password string
	}

	var body Body
	err := c.Bind(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed bind body",
		})
		return
	}

	// hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed hash password",
		})
		return
	}

	// insert to table user
	var user user.User
	user.Email = body.Email
	user.Password = string(hashPassword)
	err = Db.Create(&user).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed save user",
		})
		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func Signin(c *gin.Context) {
	// read body
	type Body struct {
		Email    string
		Password string
	}

	var body Body
	err := c.Bind(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed read body",
		})
		return
	}

	// check email
	var user user.User
	result := Db.Where(&Body{Email: body.Email}).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// create token
	key := []byte(os.Getenv("SECRET_JWT"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": user.Email,
			"exp": time.Now().Add(time.Hour * 24).Unix(), // add 24 hours
		},
	)
	tokenString, err := token.SignedString(key)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed create token",
		})
		return
	}

	// set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600*24, "", "", true, false)

	// respond with token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
