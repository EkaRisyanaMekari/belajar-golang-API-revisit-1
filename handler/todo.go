package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (handler *todoHandler) HandleUpdateTodo(c *gin.Context) {
	userSigned := c.MustGet("user").(user.User)

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	updatedTodo := todo.Todo{}

	var isExist bool
	isExist = handler.todoService.CheckTodoExistence(id)

	if !isExist {
		c.JSON(http.StatusNotFound, gin.H{
			"data": "Not found",
		})
		return
	}

	isExist = handler.todoService.CheckTodoOwnership(int(userSigned.ID), id)

	if !isExist {
		c.JSON(http.StatusForbidden, gin.H{
			"data": "Forbidden to conduct this action",
		})
		return
	}

	updatedTodo, _ = handler.todoService.GetTodoById(int(userSigned.ID), id)
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

	_, err = handler.todoService.Update(updatedTodo)
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

func (handler *todoHandler) HandleGetTodoById(c *gin.Context) {
	userSigned, _ := c.MustGet("user").(user.User)
	id, _ := strconv.Atoi(c.Param("id"))

	todo, result := handler.todoService.GetTodoById(0, id)

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

func (handler *todoHandler) HandleGetTodoBySearch(c *gin.Context) {
	userSigned, _ := c.MustGet("user").(user.User)
	keyword := c.Query("keyword")

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Please specify keyword",
		})
		return
	}

	todos := handler.todoService.GetListByKeyword(int(userSigned.ID), keyword)
	c.JSON(http.StatusOK, gin.H{
		"data": todos,
	})
}

func (handler *todoHandler) HandleDeleteTodoById(c *gin.Context) {
	userSigned := c.MustGet("user").(user.User)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var deletedTodo = todo.Todo{}

	var isExist bool
	isExist = handler.todoService.CheckTodoExistence(id)

	if !isExist {
		c.JSON(http.StatusNotFound, gin.H{
			"data": "Not found",
		})
		return
	}

	isExist = handler.todoService.CheckTodoOwnership(int(userSigned.ID), id)

	if !isExist {
		c.JSON(http.StatusForbidden, gin.H{
			"data": "Forbidden to conduct this action",
		})
		return
	}

	deletedTodo, _ = handler.todoService.GetTodoById(int(userSigned.ID), id)
	_, err = handler.todoService.Delete(deletedTodo)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"Deleted": deletedTodo,
	})
}

func (handler *todoHandler) HandleUpdateTodoStatus(c *gin.Context) {
	userSigned := c.MustGet("user").(user.User)
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	updatedTodo := todo.Todo{}

	var isExist bool
	isExist = handler.todoService.CheckTodoExistence(id)

	if !isExist {
		c.JSON(http.StatusNotFound, gin.H{
			"data": "Not found",
		})
		return
	}

	isExist = handler.todoService.CheckTodoOwnership(int(userSigned.ID), id)

	if !isExist {
		c.JSON(http.StatusForbidden, gin.H{
			"data": "Forbidden to conduct this action",
		})
		return
	}

	todoStatus := todo.TodoStatus{}
	errorBind := c.ShouldBindJSON(&todoStatus)

	if errorBind != nil {
		c.JSON(http.StatusBadRequest, errorBind)
		return
	}

	updatedTodo, _ = handler.todoService.GetTodoById(int(userSigned.ID), id)
	updatedTodo.Status = todoStatus.Status

	// errorUpdate := Db.Model(&updatedTodo).Select("Status").Updates(updatedTodo).Error
	_, errorUpdate := handler.todoService.UpdateStatus(updatedTodo)

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
