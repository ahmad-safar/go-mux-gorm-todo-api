package main

// use gorilla mux
import (
	"go-mux-gorm-todo-api/handlers"
	"go-mux-gorm-todo-api/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db := initDB()
	todo := handlers.NewTodoHandler(db)

	r := mux.NewRouter()
	r.HandleFunc("/todos", todo.GetTodos).Methods("GET")
	r.HandleFunc("/todos/{id}", todo.GetTodo).Methods("GET")
	r.HandleFunc("/todos", todo.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", todo.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}/complete", todo.CompleteTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}/uncomplete", todo.UncompleteTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", todo.DeleteTodo).Methods("DELETE")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Todo{})

	return db
}
