package handlers

import (
	"encoding/json"
	"go-mux-gorm-todo-api/models"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Data interface{}

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// getDbContext returns the database context
func getDbContext(r *http.Request) *gorm.DB {
	return r.Context().Value("db").(*gorm.DB)
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

	sendJSON(w, SuccessResponse{
		Status:  "success",
		Message: "Todos fetched successfully",
		Data:    todos,
	}, http.StatusOK)
}

// GetTodo returns a single todo
func GetTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	params := mux.Vars(r)

	db := getDbContext(r)
	result := db.First(&todo, params["id"])
	if result.Error != nil {
		sendJSONTodoNotFound(w)
		return
	}

	sendJSON(w, SuccessResponse{
		Status:  "success",
		Message: "Todo fetched successfully",
		Data:    todo,
	}, http.StatusOK)
}

// CreateTodo creates a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		sendJSON(w, ErrorResponse{
			Status: "error",
			Error:  "Invalid request body",
		}, http.StatusBadRequest)
		return
	}

	db := getDbContext(r)
	// check if todo already exists
	var existingTodo models.Todo
	db.Where("title = ?", todo.Title).First(&existingTodo)
	if existingTodo.ID != 0 {
		sendJSON(w, ErrorResponse{
			Status: "error",
			Error:  "Todo already exists",
		}, http.StatusOK)
		return
	}

	db.Create(&todo)

	sendJSON(w, SuccessResponse{
		Status:  "success",
		Message: "Todo created successfully",
		Data:    todo,
	}, http.StatusCreated)
}

// UpdateTodo updates a todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	params := mux.Vars(r)

	db := getDbContext(r)
	db.First(&todo, params["id"])
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		sendJSON(w, ErrorResponse{
			Status: "error",
			Error:  "Invalid request body",
		}, http.StatusBadRequest)
		return
	}
	db.Save(&todo)

	sendJSON(w, SuccessResponse{
		Status:  "success",
		Message: "Todo updated successfully",
		Data:    todo,
	}, http.StatusOK)
}

// toggleTodo marks a todo as completed or uncompleted
func toggleTodo(w http.ResponseWriter, r *http.Request, completed bool) {
	var todo models.Todo
	params := mux.Vars(r)

	db := getDbContext(r)
	result := db.First(&todo, params["id"])
	if result.Error != nil {
		sendJSONTodoNotFound(w)
		return
	}
	todo.Completed = completed
	db.Save(&todo)

	sendJSON(w, SuccessResponse{
		Status:  "success",
		Message: "Todo updated successfully",
		Data:    todo,
	}, http.StatusOK)
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
	result := db.First(&todo, params["id"])
	if result.Error != nil {
		sendJSONTodoNotFound(w)
		return
	}
	result = result.Delete(&todo)
	if result.Error != nil {
		sendJSONTodoNotFound(w)
	}

	sendJSON(w, SuccessResponse{
		Status:  "success",
		Message: "Todo deleted successfully",
		Data:    nil,
	}, http.StatusOK)
}

func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendJSONTodoNotFound(w http.ResponseWriter) {
	sendJSON(w, ErrorResponse{
		Status: "error",
		Error:  "Todo not found",
	}, http.StatusNotFound)
}
