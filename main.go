package main

import (
	"belajar-golang-api-revisit-1/handler"
	"belajar-golang-api-revisit-1/todo"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/belajar-golang-api-revisit-1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error DB connection")
	}
	fmt.Println("Succes DB connection")

	// DB migrattion
	db.AutoMigrate(&todo.Todo{})

	handler.Db = db

	router := gin.Default()

	v1 := router.Group("/v1")

	v1.GET("/", handler.HandleRoot)
	v1.GET("/todo/:id", handler.HandleGetTodoById)
	v1.GET("/todos", handler.HandleUrlQuery)
	v1.GET("/todos/:year/:month", handler.HandleMultiUrlParam)
	v1.POST("/todos", handler.HandlePostTodo)

	router.Run(":7878")
}
