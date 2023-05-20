package todo

type TodoInput struct {
	Title       string `binding:"required"`
	Description string `binding:"required"`
	DueDate     string `binding:"required" json:"due_date"`
	Price       int32  `binding:"required,number"`
}
