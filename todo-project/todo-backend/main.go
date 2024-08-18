package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

var (
	mutex    sync.Mutex
	db       *sqlx.DB
	nats_url string
)

const healthCheckPort = "3541"

type Todo struct {
	Id   int    `db:"id" json:"Id,omitempty"`
	Todo string `db:"todo" json:"Todo"`
	Done bool   `db:"done" json:"Done"`
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

	nats_url = os.Getenv("NATS_URL")
	if nats_url == "" {
		log.Fatal("NATS_URL is not set")
	}

	db, err = sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	createTableQuery := `CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		todo TEXT,
		Done BOOLEAN NOT NULL DEFAULT FALSE
	)`
	if _, err = db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create todo table: %v", err)
	}

	todoMux := http.NewServeMux()
	todoMux.HandleFunc("/", todosHandler)

	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/healthz", health)

	go func() {
		log.Printf("Server started on port %s", port)
		if err := http.ListenAndServe(":"+port, todoMux); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	go func() {
		if err := http.ListenAndServe(":"+healthCheckPort, healthMux); err != nil {
			log.Fatalf("Failed to start healthz endpoint: %v", err)
		}
	}()

	select {}
}

func HandleTodoPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if todo.Todo == "" {
		log.Printf("Todo cannot be empty.")
		http.Error(w, "Todo cannot be empty", http.StatusBadRequest)
		return
	}

	if !todo.Done {
		todo.Done = false
	}

	if len(todo.Todo) > 140 {
		log.Printf("Rejected: Todo exceeds 140 characters: %s", todo.Todo)
		http.Error(w, "Rejected: Todo exceeds 140 characters.", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	todoInsert := `INSERT INTO todos (todo) VALUES ($1)`
	_, err := db.Exec(todoInsert, todo.Todo)
	if err != nil {
		log.Printf("Failed to insert into database.")
		http.Error(w, "Failed to insert into database", http.StatusInternalServerError)
		return
	}

	log.Printf("Received submission: Todo=%s", todo.Todo)

	nc, err := nats.Connect(nats_url, nats.Name("API PublishBytes"))
	if err != nil {
		http.Error(w, "Failed to connect to nats", http.StatusInternalServerError)
	}

	defer nc.Close()

	msg := map[string]string{
		"user":    "bot",
		"message": "A todo was created.",
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	if err := nc.Publish("todos", msgBytes); err != nil {
		http.Error(w, "Failed to publish todos to nc", http.StatusInternalServerError)
		return
	}

	var todoList []Todo
	err = db.Select(&todoList, `SELECT todo FROM todos WHERE Done = FALSE`)
	if err != nil {
		log.Printf("Failed to fetch todos from database.")
		http.Error(w, "Failed to fetch todos from database.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		log.Printf("Failed to encode response.")
		http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
	}
}

func HandleTodoPut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if todo.Todo == "" {
		log.Printf("Todo cannot be empty.")
		http.Error(w, "Todo cannot be empty", http.StatusBadRequest)
		return
	}

	if !todo.Done {
		todo.Done = false
	}

	if len(todo.Todo) > 140 {
		log.Printf("Rejected: Todo exceeds 140 characters: %s", todo.Todo)
		http.Error(w, "Rejected: Todo exceeds 140 characters.", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	query := `UPDATE todos SET todo = $1, done = $2 WHERE id = $3`
	_, err := db.Exec(query, todo.Todo, todo.Done, todo.Id)
	if err != nil {
		log.Printf("Failed to update the todo with ID %s, %v", todo.Id, err)
		http.Error(w, "Failed to update the todo in the database", http.StatusInternalServerError)
		return
	}

	nc, err := nats.Connect(nats_url, nats.Name("API PublishBytes"))
	if err != nil {
		http.Error(w, "Failed to connect to nats", http.StatusInternalServerError)
	}

	msg := map[string]string{
		"user":    "bot",
		"message": "A todo was updated.",
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	if err := nc.Publish("todos", msgBytes); err != nil {
		http.Error(w, "Failed to publish to nats", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleTodoGet(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")

	var todoList []Todo
	err := db.Select(&todoList, `SELECT * FROM todos WHERE done=false`)
	if err != nil {
		log.Printf("Failed to fetch todos from database.")
		http.Error(w, "Failed to fetch todos from database.", http.StatusInternalServerError)
		todoList = []Todo{}
	}

	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		log.Printf("Failed to encode response.")
		http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
	}
}

func HandleTodoDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("Id cannot be empty.")
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	_, err := db.Exec(`DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		log.Printf("Failed to delete todo with id %s from the database: %v", id, err)
		http.Error(w, "Failed to delete todo from database.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandleTodoPost(w, r)
	case http.MethodGet:
		HandleTodoGet(w, r)
	case http.MethodPut:
		HandleTodoPut(w, r)
	case http.MethodDelete:
		HandleTodoDelete(w, r)
	case http.MethodOptions:
		w.Header().Set("Allow", "POST, GET, PUT, DELETE")
		w.WriteHeader(http.StatusOK)
	default:
		w.Header().Set("Allow", "POST, GET, PUT, DELETE")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	err := db.Ping()
	if err != nil {
		http.Error(w, "Database connection failed.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
