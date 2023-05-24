package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"belajar-golang-api-revisit-1/todo"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var Db *gorm.DB

func HandleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    "Eka",
		"address": "Bandung x",
	})
}

func HandleGetTodoById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func HandleMultiUrlParam(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	c.JSON(http.StatusOK, gin.H{
		"year":  year,
		"month": month,
	})
}

func HandlePostTodo(c *gin.Context) {
	var todoInput todo.TodoInput

	err := c.ShouldBindJSON(&todoInput)
	if err != nil {
		var jsonError *json.UnmarshalTypeError
		fmt.Print(err)
		errMessages := []string{}
		if errors.As(err, &jsonError) {
			msg := fmt.Sprintf("Error [field: %s], actual type: %s", err.(*json.UnmarshalTypeError).Field, err.(*json.UnmarshalTypeError).Value)
			errMessages = append(errMessages, msg)
		} else {
			for _, e := range err.(validator.ValidationErrors) {
				msg := fmt.Sprintf("Error [field: %s], is: %s", e.Field(), e.ActualTag())
				errMessages = append(errMessages, msg)
			}
		}
		c.JSON(http.StatusBadRequest, errMessages)
		return
	}

	newTodo := todo.Todo{}
	newTodo.Title = todoInput.Title
	newTodo.Description = todoInput.Description
	newTodo.DueDate = todoInput.DueDate
	newTodo.Status = 0

	err = Db.Create(&newTodo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		fmt.Println(err)
		log.Fatal("Error create record")
	}

	c.JSON(http.StatusOK, newTodo)

}

func HandleGetTodosByStatus(c *gin.Context) {
	status := c.Query("status")

	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Please specify status",
		})
		return
	}

	var todos []todo.Todo
	Db.Find(&todos, "status = ?", status)
	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})

}
