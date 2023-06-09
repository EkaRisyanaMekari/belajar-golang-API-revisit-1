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
