package models

type Todos struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	AddedAt   int64  `json:"added_at"`
	UpdatedAt int64  `json:"updated_at"`
	Completed bool   `json:"completed"`
}

type CreateTodoInput struct {
	Title string `json:"title" binding:"required"`
}

type UpdateTodoInput struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
