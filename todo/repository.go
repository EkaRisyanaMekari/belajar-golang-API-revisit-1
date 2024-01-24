package todo

import (
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	Create(todo Todo) (Todo, error)
	GetListAll(userId int) []Todo
	GetListByStatus(userId int, status string) []Todo
	GetFirst(userId int, id int) (Todo, *gorm.DB)
	GetListByKeyword(userId int, keyword string) []Todo
	Delete(todo Todo) (Todo, error)
	// Update(todo TodoInput) (Todo, error)
	// Delete() (Todo, error)
	// GetList() ([]Todo, error)
	// GetDetail() (Todo, error)
	// UpdateStatus() (Todo, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(todo Todo) (Todo, error) {
	error := r.db.Create(&todo).Error
	return todo, error
}

func (r *repository) GetListAll(userId int) []Todo {
	var todos []Todo
	r.db.Where(&Todo{UserId: userId}).Find(&todos)
	return todos
}

func (r *repository) GetListByStatus(userId int, status string) []Todo {
	var todos []Todo
	r.db.Where(&Todo{UserId: userId}).Find(&todos, "status = ?", status)
	return todos
}

func (r *repository) GetFirst(userId int, id int) (Todo, *gorm.DB) {
	var todo Todo
	var result *gorm.DB
	if userId == 0 {
		result = r.db.First(&todo, id)
	} else {
		result = r.db.First(&todo, Todo{ID: id, UserId: userId})
	}
	return todo, result
}

func (r *repository) GetListByKeyword(userId int, keyword string) []Todo {
	var todos []Todo
	conditions := []string{"%", keyword, "%"}
	r.db.Where("description like ?", strings.Join(conditions, "")).Where(&Todo{UserId: userId}).Find(&todos)
	return todos
}

func (r *repository) Delete(todo Todo) (Todo, error) {
	error := r.db.Delete(&todo).Error
	return todo, error
}
