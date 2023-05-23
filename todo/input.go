package todo

import "time"

type TodoInput struct {
	Title       string    `binding:"required"`
	Description string    `binding:"required"`
	DueDate     time.Time `binding:"required" json:"due_date"`
	Price       int32     `binding:"required,number"`
}
