package todo

import "gorm.io/gorm"

type Service interface {
	Create(todo Todo) (Todo, error)
	GetListByStatus(userId int, status string) []Todo
	GetTodoById(userId int, id int) (Todo, *gorm.DB)
	GetListByKeyword(userId int, keyword string) []Todo
	Delete(todo Todo) (Todo, error)
	CheckTodoExistence(id int) bool
	CheckTodoOwnership(userId int, id int) bool
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

func (s *service) GetTodoById(userId int, id int) (Todo, *gorm.DB) {
	todo, result := s.repository.GetFirst(userId, id)
	return todo, result
}

func (s *service) GetListByKeyword(userId int, keyword string) []Todo {
	todos := s.repository.GetListByKeyword(userId, keyword)
	return todos
}

func (s *service) Delete(todo Todo) (Todo, error) {
	var deletedTodo Todo
	deletedTodo, err := s.repository.Delete(todo)
	return deletedTodo, err
}

func (s *service) CheckTodoExistence(id int) bool {
	_, result := s.repository.GetFirst(0, id)
	return result.RowsAffected > 0
}

func (s *service) CheckTodoOwnership(userId int, id int) bool {
	_, result := s.repository.GetFirst(userId, id)
	return result.RowsAffected > 0
}
