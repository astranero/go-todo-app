package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	mutex sync.Mutex
	db    *sqlx.DB
)

type Todo struct {
	Todo string `db:"todo"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file.")
	}

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err = sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS todos (
		todo TEXT
	)`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create todo table: %v", err)
	}

	http.HandleFunc("/todos", todosHandler)

	log.Printf("Server started on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func HandleTodoPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	todo := r.FormValue("todo")
	if todo == "" {
		http.Error(w, "Todo cannot be empty", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	var todoList []Todo
	err := db.Select(&todoList, `SELECT todo FROM todos`)
	if err != nil {
		http.Error(w, "Failed to fetch todos from database.", http.StatusInternalServerError)
		return
	}

	log.Printf("Received submission: Todo=%s", todo)
	todoInsert := `INSERT INTO todos (todo) VALUES ($1)`
	_, err = db.Exec(todoInsert, todo)
	if err != nil {
		http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
		return
	}

	todoList = append(todoList, Todo{Todo: todo})
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func HandleTodoGet(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")

	var todoList []Todo
	err := db.Select(&todoList, `SELECT todo FROM todos`)
	if err != nil {
		http.Error(w, "Failed to fetch todos from database.", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandleTodoPost(w, r)
	case http.MethodGet:
		HandleTodoGet(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
