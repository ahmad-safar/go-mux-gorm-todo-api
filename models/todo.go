package models

import (
	"time"

	"gorm.io/gorm"
)

// gorm.Model definition
type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Todo is a model for a todo item
type Todo struct {
	Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Todos is a slice of Todo
type Todos []Todo
