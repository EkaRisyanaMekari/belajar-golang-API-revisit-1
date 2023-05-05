package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", handleRoot)
	router.GET("/todo/:id", handleGetTodoById)
	router.GET("/todos", handleUrlQuery)
	router.GET("/todos/:year/:month", handleMultiUrlParam)

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
	c.JSON(http.StatusOK, gin.H{
		"title": title,
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
