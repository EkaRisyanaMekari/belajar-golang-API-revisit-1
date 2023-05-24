package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func HandleGetTodoById(c *gin.Context) {
	id := c.Param("id")

	var todo todo.Todo
	result := Db.First(&todo, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "Not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": todo,
	})
}

func HandleGetTodoBySearch(c *gin.Context) {
	keyword := c.Query("keyword")

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Please specify keyword",
		})
		return
	}

	var todos []todo.Todo
	conditions := []string{"%", keyword, "%"}
	Db.Where("description like ?", strings.Join(conditions, "")).Find(&todos)
	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}
