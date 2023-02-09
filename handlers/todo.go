package handlers

import (
	"encoding/json"
	"errors"
	"go-mux-gorm-todo-api/models"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

// GetTodos returns all todos
func (t *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos models.Todos

	// get type of todos to return
	todoType := r.URL.Query().Get("type")
	if todoType == "" || (todoType != "completed" && todoType != "uncompleted") {
		todoType = "all"
	}

	switch todoType {
	case "all":
		t.db.Find(&todos)
	case "completed":
		t.db.Where("completed = ?", true).Find(&todos)
	case "uncompleted":
		t.db.Where("completed = ?", false).Find(&todos)
	}

	respondWithJSON(w, http.StatusOK, "Todos fetched successfully", &todos)
}

// GetTodo returns a single todo
func (t *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	params := mux.Vars(r)

	if err := t.db.First(&todo, params["id"]).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Todo not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, "Todo fetched successfully", &todo)
}

// CreateTodo creates a new todo
func (t *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	var existingTodo models.Todo
	t.db.Where("title = ?", todo.Title).First(&existingTodo)
	if existingTodo.ID != 0 {
		respondWithError(w, http.StatusBadRequest, "Todo already exists", errors.New("todo already exists"))
		return
	}

	if err := t.db.Create(&todo).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Todo could not be created", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, "Todo created successfully", &todo)
}

// UpdateTodo updates a todo
func (t *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	params := mux.Vars(r)

	if err := t.db.First(&todo, params["id"]).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Todo not found", err)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if err := t.db.Save(&todo).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Todo could not be updated", err)
		return
	}

	respondWithJSON(w, http.StatusOK, "Todo updated successfully", &todo)
}

// toggleTodo marks a todo as completed or uncompleted
func (t *TodoHandler) toggleTodo(w http.ResponseWriter, r *http.Request, completed bool) {
	var todo models.Todo
	params := mux.Vars(r)

	if err := t.db.First(&todo, params["id"]).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Todo not found", err)
		return
	}
	todo.Completed = completed
	if err := t.db.Save(&todo).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Todo could not be updated", err)
		return
	}

	respondWithJSON(w, http.StatusOK, "Todo updated successfully", &todo)
}

// CompleteTodo marks a todo as completed
func (t *TodoHandler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	t.toggleTodo(w, r, true)
}

// UncompletedTodo marks a todo as uncompleted
func (t *TodoHandler) UncompleteTodo(w http.ResponseWriter, r *http.Request) {
	t.toggleTodo(w, r, false)
}

// DeleteTodo deletes a todo
func (t *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	params := mux.Vars(r)

	// check if todo exists
	res := t.db.First(&todo, params["id"])
	if res.Error != nil {
		respondWithError(w, http.StatusNotFound, "Todo not found", res.Error)
		return
	}
	if err := res.Delete(&todo).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Todo could not be deleted", err)
		return
	}

	respondWithJSON(w, http.StatusOK, "Todo deleted successfully", &todo)
}
