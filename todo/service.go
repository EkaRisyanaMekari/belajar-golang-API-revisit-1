package todo

import "gorm.io/gorm"

type Service interface {
	Create(todo Todo) (Todo, error)
	GetListByStatus(userId int, status string) []Todo
	GetTodoById(userId int, id string) (Todo, *gorm.DB)
	// Update(todo TodoInput) (Todo, error)
	// Delete() (Todo, error)
	// GetList() ([]Todo, error)
	// GetDetail() (Todo, error)
	// UpdateStatus() (Todo, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) Create(todo Todo) (Todo, error) {
	newTodo, err := s.repository.Create(todo)
	return newTodo, err
}

func (s *service) GetListByStatus(userId int, status string) []Todo {
	var todos []Todo
	if status == "" {
		todos = s.repository.GetListAll(userId)
	} else {
		todos = s.repository.GetListByStatus(userId, status)
	}
	return todos
}

func (s *service) GetTodoById(userId int, id string) (Todo, *gorm.DB) {
	todo, result := s.repository.GetFirst(userId, id)
	return todo, result
}
