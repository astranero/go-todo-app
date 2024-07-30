package main

import (
	"net/http"
	"log"
	"os"
	"sync"
	"encoding/json"
	"github.com/joho/godotenv"
)

var (
	mutex sync.Mutex
	todoList []string
)

func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load .env file.")
	}
	
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
	}


	http.HandleFunc("/todos", todosHandler)

	log.Printf("Server started in port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


func HandleTodoPost(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	todo := r.FormValue("todo")
	if todo == "" {
		http.Error(w, "Todo cannot be empty", http.StatusBadRequest)
		return 
	}

	log.Printf("Received submission: Todo=%s", todo)
	mutex.Lock()
	todoList = append(todoList, todo)
	mutex.Unlock()

	w.Header().Set("Content-Type", "text/html")
	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func HandleTodoGet(w http.ResponseWriter, r *http.Request){
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func todosHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodPost:
		HandleTodoPost(w, r)
	case http.MethodGet:
		HandleTodoGet(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}