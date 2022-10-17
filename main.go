package main

// use gorilla mux
import (
	"context"
	"go-mux-gorm-todo-api/constants"
	"go-mux-gorm-todo-api/handlers"
	"go-mux-gorm-todo-api/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	DB = initDB()

	r := mux.NewRouter()
	r.Use(dbMiddleware)
	r.HandleFunc("/todos", handlers.GetTodos).Methods("GET")
	r.HandleFunc("/todos/{id}", handlers.GetTodo).Methods("GET")
	r.HandleFunc("/todos", handlers.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}/complete", handlers.CompleteTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}/uncomplete", handlers.UncompleteTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods("DELETE")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func dbMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), constants.DbKey{}, DB)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Todo{})

	return db
}
