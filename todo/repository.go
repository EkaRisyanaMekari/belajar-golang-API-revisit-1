package todo

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(todo Todo) (Todo, error)
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
