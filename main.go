package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")

	v1.GET("/", handleRoot)
	v1.GET("/todo/:id", handleGetTodoById)
	v1.GET("/todos", handleUrlQuery)
	v1.GET("/todos/:year/:month", handleMultiUrlParam)
	v1.POST("/todos", handlePostTodo)

	router.Run(":7878")
}

func handleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    "Eka",
		"address": "Bandung x",
	})
}

func handleGetTodoById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func handleUrlQuery(c *gin.Context) {
	title := c.Query("title")
	price := c.Query("price")
	c.JSON(http.StatusOK, gin.H{
		"title": title,
		"price": price,
	})
}

func handleMultiUrlParam(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	c.JSON(http.StatusOK, gin.H{
		"year":  year,
		"month": month,
	})
}

type TodoInput struct {
	Title       string `binding:"required"`
	Description string `binding:"required"`
	DueDate     string `binding:"required" json:"due_date"`
	Price       int32  `binding:"required,number"`
}

func handlePostTodo(c *gin.Context) {
	var todoInput TodoInput

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
