package todo

import "time"

type Todo struct {
	ID          int
	Title       string
	Description string
	DueDate     time.Time
	Status      int8
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
	UserId      int
}
