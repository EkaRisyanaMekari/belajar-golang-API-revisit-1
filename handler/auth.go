package handler

import (
	"belajar-golang-api-revisit-1/user"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// set cookie

	// respond with token
	c.JSON(http.StatusOK, gin.H{
		"token": "success signin",
	})
}
