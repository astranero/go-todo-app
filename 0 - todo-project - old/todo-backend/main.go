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

	go func(){
		http.HandleFunc("/", todosHandler)

		log.Printf("Server started on port %s", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func(){
		http.HandleFunc("/healthz", health)
		port := "3541"
		http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Failed to start healthz endpoint: %v", err)
		}
	}()
}

func HandleTodoPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	todo := r.FormValue("todo")
	if todo == "" {
		log.Printf("Todo cannot be empty.")
		http.Error(w, "Todo cannot be empty", http.StatusBadRequest)
		return
	}

	if len(todo) > 140 {
        log.Printf("Rejected: Todo exceeds 140 characters: %s", todo)
        http.Error(w, "Rejected: Todo exceeds 140 characters.", http.StatusBadRequest)
        return
    }

	mutex.Lock()
	defer mutex.Unlock()

	todoInsert := `INSERT INTO todos (todo) VALUES ($1)`
	_, err := db.Exec(todoInsert, todo)
	if err != nil {
		log.Printf("Failed to insert into database.")
		http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
		return
	}

	log.Printf("Received submission: Todo=%s", todo)

	var todoList []Todo
	err = db.Select(&todoList, `SELECT todo FROM todos`)
	if err != nil {
		log.Printf("Failed to fetch todos from database.")
		http.Error(w, "Failed to fetch todos from database.", http.StatusInternalServerError)
		return
	}

	todoList = append(todoList, Todo{Todo: todo})
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		log.Printf("Failed to encode response.")
		http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
	}
}

func HandleTodoGet(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")

	var todoList []Todo
	err := db.Select(&todoList, `SELECT todo FROM todos`)
	if err != nil {
		log.Printf("Failed to fetch todos from database.")
		http.Error(w, "Failed to fetch todos from database.", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		log.Printf("Failed to encode response.")
		http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
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

func health(w http.ResponseWriter, r *http.Request){
	err := db.Ping()
	if err != nil {
		http.Error(w, "Database connection failed.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w,"OK")
}
