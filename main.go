package main

import (
	"belajar-golang-api-revisit-1/handler"
	"belajar-golang-api-revisit-1/middleware"
	"belajar-golang-api-revisit-1/todo"
	"belajar-golang-api-revisit-1/user"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// only load the .env file when running locally
	// check for a RAILWAY_ENVIRONMENT, if not found, code is running locally
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); exists == false {
	    if err := godotenv.Load(); err != nil {
	        log.Fatal("error loading .env file:", err)
	    }
	}
	
	dsn := os.Getenv("DATABASE")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error DB connection", err)
	}
	fmt.Println("Succes DB connection")

	// DB migrattion
	db.AutoMigrate(&todo.Todo{})
	db.AutoMigrate(&user.User{})

	handler.Db = db

	todoRepository := todo.NewRepository(db)
	todoService := todo.NewService(todoRepository)
	todoHandler := handler.NewTodoHandler(todoService)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://learn-typescript-vue3.vercel.app/"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	v1 := router.Group("/v1")

	router.GET("/", handler.HandleRoot)
	v1.POST("/signup", handler.Signup)
	v1.POST("/signin", handler.Signin)
	v1.POST("/todos", middleware.RequireAuth, todoHandler.HandlePostTodo)
	v1.PUT("/todos", middleware.RequireAuth, todoHandler.HandleUpdateTodo)
	v1.PUT("/todos/update-status", middleware.RequireAuth, todoHandler.HandleUpdateTodoStatus)
	v1.GET("/todos", middleware.RequireAuth, todoHandler.HandleGetTodosByStatus)
	v1.GET("/todos/:id", middleware.RequireAuth, todoHandler.HandleGetTodoById)
	v1.DELETE("/todos/:id", middleware.RequireAuth, todoHandler.HandleDeleteTodoById)
	v1.GET("/todos/search", middleware.RequireAuth, todoHandler.HandleGetTodoBySearch)

	router.Run(":" + os.Getenv("PORT"))
}
