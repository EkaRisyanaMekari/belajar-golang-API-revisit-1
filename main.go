package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", handleRoot)
	router.GET("/todo/:id", handleGetTodoById)
	router.GET("/todos", handleUrlQuery)
	router.GET("/todos/:year/:month", handleMultiUrlParam)
	router.POST("/todos", handlePostTodo)

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
	Title       string
	Description string
}

func handlePostTodo(c *gin.Context) {
	var todoInput TodoInput

	err := c.ShouldBindJSON(&todoInput)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"title": todoInput.Title,
		"desc":  todoInput.Description,
	})

}
