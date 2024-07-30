package main

import (
	"net/http"
	"log"
	"os"
	"time"
	"sync"
	"io"
	"fmt"
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
	
	port := os.Getenv("GO_PORT")
	if port == "" {
		port = "8080"
	}

	imagePath := os.Getenv("imagePath")
	imageUrl := os.Getenv("imageUrl")

	err = downloadFile(imagePath, imageUrl)
	if err != nil {
		log.Fatalf("Failed to download image.")
	}

	ticker := time.NewTicker(3600 * time.Second)
	go func(){
		for range ticker.C {
			mutex.Lock()
			err := downloadFile(imagePath, imageUrl)
			if err != nil {
				log.Fatalf("Failed to download image.")
			}
			mutex.Unlock()
		}
	}()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)	

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w,r,imagePath)
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request){
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
		w.Header().Set("Content-Type", "text/html")
		mutex.Lock()
		todoList = append(todoList, todo)
		mutex.Unlock()
		json.NewEncoder(w).Encode(todoList)
	})

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request){
		mutex.Lock()
		defer mutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todoList)
	})

	log.Printf("Server started in port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func downloadFile(filepath string, url string)(err error){
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}