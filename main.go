package main

import (
	"belajar-golang-api-revisit-1/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")

	v1.GET("/", handler.HandleRoot)
	v1.GET("/todo/:id", handler.HandleGetTodoById)
	v1.GET("/todos", handler.HandleUrlQuery)
	v1.GET("/todos/:year/:month", handler.HandleMultiUrlParam)
	v1.POST("/todos", handler.HandlePostTodo)

	router.Run(":7878")
}
