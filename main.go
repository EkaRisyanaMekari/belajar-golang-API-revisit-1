package main

import (
	"belajar-golang-api-revisit-1/handler"
	"belajar-golang-api-revisit-1/middleware"
	"belajar-golang-api-revisit-1/todo"
	"belajar-golang-api-revisit-1/user"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error load env")
		fmt.Println(err)
	}
	dsn := os.Getenv("DATABASE")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error DB connection")
	}
	fmt.Println("Succes DB connection")

	// DB migrattion
	db.AutoMigrate(&todo.Todo{})
	db.AutoMigrate(&user.User{})

	handler.Db = db

	router := gin.Default()

	v1 := router.Group("/v1")

	v1.GET("/", middleware.RequireAuth, handler.HandleRoot)
	v1.POST("/signup", handler.Signup)
	v1.POST("/signin", handler.Signin)
	v1.POST("/todos", middleware.RequireAuth, handler.HandlePostTodo)
	v1.PUT("/todos", middleware.RequireAuth, handler.HandleUpdateTodo)
	v1.GET("/todos", middleware.RequireAuth, handler.HandleGetTodosByStatus)
	v1.GET("/todos/:id", middleware.RequireAuth, handler.HandleGetTodoById)
	v1.DELETE("/todos/:id", middleware.RequireAuth, handler.HandleDeleteTodoById)
	v1.GET("/todos/search", middleware.RequireAuth, handler.HandleGetTodoBySearch)

	router.Run(":7878")
}
