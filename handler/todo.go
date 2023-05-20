package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"belajar-golang-api-revisit-1/todo"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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

func HandleUrlQuery(c *gin.Context) {
	title := c.Query("title")
	price := c.Query("price")
	c.JSON(http.StatusOK, gin.H{
		"title": title,
		"price": price,
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

	c.JSON(http.StatusOK, gin.H{
		"title":    todoInput.Title,
		"desc":     todoInput.Description,
		"due_date": todoInput.DueDate,
		"price":    todoInput.Price,
	})

}
