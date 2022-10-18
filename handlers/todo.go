package handlers

import (
	"encoding/json"
	"go-mux-gorm-todo-api/constants"
	"go-mux-gorm-todo-api/models"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// getDbContext returns the database context
func getDbContext(r *http.Request) *gorm.DB {
	return r.Context().Value(constants.DbKey{}).(*gorm.DB)
}

// GetTodos returns all todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos models.Todos

	// get type of todos to return
	todoType := r.URL.Query().Get("type")
	if todoType == "" || (todoType != "completed" && todoType != "uncompleted") {
		todoType = "all"
	}

	db := getDbContext(r)
	switch todoType {
	case "all":
		db.Find(&todos)
	case "completed":
		db.Where("completed = ?", true).Find(&todos)
	case "uncompleted":
		db.Where("completed = ?", false).Find(&todos)
	}

	respondWithJSON(w, http.StatusOK, "Todos fetched successfully", &todos)
}

// GetTodo returns a single todo
func GetTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	params := mux.Vars(r)

	db := getDbContext(r)
	if err := db.First(&todo, params["id"]).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}

	respondWithJSON(w, http.StatusOK, "Todo fetched successfully", &todo)
}

// CreateTodo creates a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	db := getDbContext(r)
	// check if todo already exists
	var existingTodo models.Todo
	db.Where("title = ?", todo.Title).First(&existingTodo)
	if existingTodo.ID != 0 {
		respondWithError(w, http.StatusBadRequest, "Todo already exists")
		return
	}

	if err := db.Create(&todo).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusCreated, "Todo created successfully", &todo)
}

// UpdateTodo updates a todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	params := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	db := getDbContext(r)
	if err := db.First(&todo, params["id"]).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}
	if err := db.Save(&todo).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, "Todo updated successfully", &todo)
}

// toggleTodo marks a todo as completed or uncompleted
func toggleTodo(w http.ResponseWriter, r *http.Request, completed bool) {
	var todo models.Todo
	params := mux.Vars(r)

	db := getDbContext(r)
	if err := db.First(&todo, params["id"]).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}
	todo.Completed = completed
	if err := db.Save(&todo).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, "Todo updated successfully", &todo)
}

// CompleteTodo marks a todo as completed
func CompleteTodo(w http.ResponseWriter, r *http.Request) {
	toggleTodo(w, r, true)
}

// UncompletedTodo marks a todo as uncompleted
func UncompleteTodo(w http.ResponseWriter, r *http.Request) {
	toggleTodo(w, r, false)
}

// DeleteTodo deletes a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	params := mux.Vars(r)

	db := getDbContext(r)
	// check if todo exists
	res := db.First(&todo, params["id"])
	if res.Error != nil {
		respondWithError(w, http.StatusNotFound, "Todo not found")
		return
	}
	if err := res.Delete(&todo).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, "Todo deleted successfully", nil)
}
