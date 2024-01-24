package todo

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(todo Todo) (Todo, error)
	GetListAll(userId int) []Todo
	GetListByStatus(userId int, status string) []Todo
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
