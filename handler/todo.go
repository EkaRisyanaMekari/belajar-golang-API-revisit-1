package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"belajar-golang-api-revisit-1/todo"
	"belajar-golang-api-revisit-1/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var Db *gorm.DB

type todoHandler struct {
	todoService todo.Service
}

func NewTodoHandler(todoService todo.Service) *todoHandler {
	return &todoHandler{todoService}
}

func HandleRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    "Hello World",
		"address": "Bandung",
	})
}

func (handler *todoHandler) HandlePostTodo(c *gin.Context) {
	userSigned := c.MustGet("user").(user.User)
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
	newTodo.UserId = int(userSigned.ID)

	createdTodo, err := handler.todoService.Create(newTodo)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		fmt.Println(err)
		log.Fatal("Error create record")
	}

	c.JSON(http.StatusOK, createdTodo)

}

func HandleUpdateTodo(c *gin.Context) {
	userSigned := c.MustGet("user").(user.User)

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	updatedTodo := todo.Todo{}

	result := Db.Where(&todo.Todo{ID: int(id)}).Find(&updatedTodo)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "Not found",
		})
		return
	}

	result = Db.Where(&todo.Todo{ID: int(id), UserId: int(userSigned.ID)}).Find(&updatedTodo)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"data": "Forbidden to conduct this action",
		})
		return
	}

	var todoInput todo.TodoInput

	err = c.ShouldBindJSON(&todoInput)
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

	updatedTodo.Title = todoInput.Title
	updatedTodo.Description = todoInput.Description
	updatedTodo.DueDate = todoInput.DueDate
	updatedTodo.ID = id

	err = Db.Model(&updatedTodo).Select("Title", "Description", "DueDate").Updates(updatedTodo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		fmt.Println(err)
		log.Fatal("Error update record")
	}

	c.JSON(http.StatusOK, updatedTodo)

}

func (handler *todoHandler) HandleGetTodosByStatus(c *gin.Context) {
	userSigned, _ := c.MustGet("user").(user.User)
	status := c.Query("status")

	todos := handler.todoService.GetListByStatus(int(userSigned.ID), status)
	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})

}

func HandleGetTodoById(c *gin.Context) {
	userSigned, _ := c.MustGet("user").(user.User)
	id := c.Param("id")

	var todo todo.Todo
	result := Db.First(&todo, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "Not found",
		})
		return
	}

	if todo.UserId != int(userSigned.ID) {
		c.JSON(http.StatusForbidden, gin.H{
			"data": "Forbidden to get this data",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": todo,
	})
}

func HandleGetTodoBySearch(c *gin.Context) {
	userSigned, _ := c.MustGet("user").(user.User)
	keyword := c.Query("keyword")

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Please specify keyword",
		})
		return
	}

	var todos []todo.Todo
	conditions := []string{"%", keyword, "%"}
	Db.Where("description like ?", strings.Join(conditions, "")).Where(&todo.Todo{UserId: int(userSigned.ID)}).Find(&todos)
	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}

func HandleDeleteTodoById(c *gin.Context) {
	userSigned := c.MustGet("user").(user.User)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var deletedTodo = todo.Todo{}

	result := Db.Where(&todo.Todo{ID: int(id)}).Find(&deletedTodo)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"data": "Not found",
		})
		return
	}

	result = Db.Where(&todo.Todo{ID: int(id), UserId: int(userSigned.ID)}).Find(&deletedTodo)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"data": "Forbidden to conduct this action",
		})
		return
	}

	err = Db.Delete(&deletedTodo).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"Deleted": deletedTodo,
	})
}

func HandleUpdateTodoStatus(c *gin.Context) {
	userSigned := c.MustGet("user").(user.User)
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	updatedTodo := todo.Todo{}

	result := Db.Where(&todo.Todo{ID: int(id), UserId: int(userSigned.ID)}).Find(&updatedTodo)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"data": "No data updated",
		})
		return
	}

	todoStatus := todo.TodoStatus{}
	errorBind := c.ShouldBindJSON(&todoStatus)

	if errorBind != nil {
		c.JSON(http.StatusBadRequest, errorBind)
		return
	}

	updatedTodo.Status = todoStatus.Status

	errorUpdate := Db.Model(&updatedTodo).Select("Status").Updates(updatedTodo).Error

	if errorUpdate != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "Error when update status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "Success update status",
	})

}
